import base64
import socket
import requests

from typing import Union

from flask import (
    json,
    render_template,
    request,
    redirect,
    session,
    url_for,
    flash,
    abort,
    current_app,
)
from flask_login import login_user, logout_user, login_required
from flask_mail import Message
from werkzeug.wrappers import response

from submission import mail
from submission.auth import bp_auth
from submission.models import User, Token
from submission.models.users import UserRole
from .forms.login import LoginForm, UserEmailForm, PasswordResetForm


@bp_auth.route("/login", methods=["GET"])
def login() -> str:
    """Renders login form"""
    form = LoginForm(request.form)
    return render_template("auth/login.html.j2", form=form)


@bp_auth.route("/login", methods=["POST"])
def login_post() -> response.Response:
    """Handles POST requests to the login page to log in users"""
    email = request.form.get("username")
    password = request.form.get("password")
    remember = True if request.form.get("remember") else False

    if remember:
        session.permanent = True

    response = requests.post(
        f"{current_app.config['API_BASE']}/login",
        json={"email": email, "password": password},
    )

    if response.status_code != 200:
        flash("Please check your login details", "warning")
        return redirect(url_for("auth.login"))

    token_string = response.json().get("token")

    session["token"] = token_string

    if not token_string:
        flash("Please check your login details", "warning")
        return redirect(url_for("auth.login"))

    token_parts = token_string.split(".")
    if len(token_parts) != 3:
        flash("Invalid token", "warning")
        return redirect(url_for("auth.login"))

    token_data = base64.urlsafe_b64decode(token_parts[1] + "==")
    if not token_data:
        flash("Invalid token", "warning")
        return redirect(url_for("auth.login"))

    try:
        token_data = json.loads(token_data)
    except json.JSONDecodeError:
        flash("Invalid token", "warning")
        return redirect(url_for("auth.login"))

    if not token_data:
        flash("Invalid token", "warning")
        return redirect(url_for("auth.login"))

    current_app.logger.debug(token_data)

    user = User.from_json(token_data["user"])

    login_user(user, remember=remember)

    return redirect(url_for("main.index"))


@bp_auth.route("/logout")
@login_required
def logout() -> response.Response:
    """Logs out current user and redirects to the login page"""
    session['token'] = None
    logout_user()
    return redirect(url_for("auth.login"))


@bp_auth.route("/reset-my-password", methods=["GET", "POST"])
def password_email() -> Union[str, response.Response]:
    """Send an email with password reset link to user provided email"""
    form = UserEmailForm(request.form)
    if request.method == "POST" and form.validate():
        email = form.email.data
        
        response = requests.post(
            f"{current_app.config['API_BASE']}/user/password/reset",
            json={"email": email},
        )

        if response.status_code != 200:
            flash("could not request password reset email. Please try again later", "error")
            return render_template("auth/pw_reset_request.html", form=form)
        
        data = response.json()

        socket.setdefaulttimeout(5)

        # TODO: send email
        mail.send(
            Message(
                subject="Change your MIBiG password",
                recipients=[email],
                body=f"Hello, click this link {current_app.config['BASE_URL']}/auth/reset/{data['email']}/{data['challenge']}",
            )
        )
        flash("Please check your email")
        return redirect(url_for("auth.login"))

    return render_template("auth/pw_reset_request.html", form=form)


@bp_auth.route("/reset/<email>/<token_id>", methods=["GET", "POST"])
def reset_password(email: str, token_id: str) -> Union[str, response.Response]:
    """Allow a user to change their password via email provided link

    Arguments:
        token_id (str): uuid token
    """

    form = PasswordResetForm(request.form)

    form['email'].data = email
    form['challenge'].data = token_id

    if request.method == "POST" and form.validate():
        response = requests.post(
            f"{current_app.config['API_BASE']}/user/password/challenge",
            json=form.data,
        )

        if response.status_code != 200:
            flash("could not reset password", "error")
            return render_template("auth/pw_reset.html", form=form)
        
        flash("Password reset successfully")
        

        flash("Successfully changed password")
        return redirect(url_for("auth.login"))

    return render_template("auth/pw_reset.html", form=form)
