from pathlib import Path

from flask import (
    current_app,
    render_template,
    request,
    redirect,
    session,
    url_for,
    flash,
)
from flask_login import login_required, current_user, login_user
import requests

from submission.main import bp_main
from submission.main.forms import SelectExisting, UserDetailsEditForm
from submission.auth import auth_role
from submission.models.users import User, UserInfo
from submission.utils import Storage


@bp_main.route("/", methods=["GET", "POST"])
@login_required
def index():
    """Main page, edit existing entry or start new entry"""
    form = SelectExisting(request.form)
    if request.method == "POST":
        # create new entry
        if form.submit.data:
            bgc_id = "new"
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))

        # edit valid existing entry
        if request.form.get("edit") and form.validate():
            bgc_id = form.accession.data
            Storage.create_entry_if_not_exists(bgc_id)

            return redirect(url_for("edit.edit_bgc"))

    return render_template("main/index.html", form=form, user_id=current_user.id)


@bp_main.route("/delete", methods=["DELETE"])
def delete() -> str:
    """Dummy route to delete any target

    Returns:
        str: Dummy value
    """
    return ""


@bp_main.route("/profile", methods=["GET", "POST"])
@login_required
def profile():
    if request.method == "POST":
        # current_user.info["name"] = request.form["name"]
        # current_user.info["call_name"] = request.form["call_name"]
        # current_user.info["orcid"] = request.form.get("orc_id", None)
        # current_user.info["organization1"] = request.form["organization1"]
        # current_user.info["organization2"] = request.form.get("organization2", None)
        # current_user.info["organization3"] = request.form.get("organization3", None)

        response = requests.patch(
            f"{current_app.config['API_BASE']}/user/{current_user.id}",
            data=current_user.to_json(),
            headers={"Authorization": "Bearer " + session["token"]},
        )

        if response.status_code != 200:
            flash("Error updating your user details: " + str(response.json()))
            return render_template("main/profile.html.j2")

        login_user(User.get_user(current_user.id))

        flash("Updated your user details")

    form = UserDetailsEditForm()

    info: UserInfo = current_user.info

    form.name.data = info.name
    form.call_name.data = info.call_name
    form.orcid.data = info.orc_id
    form.organisation.data = info.organisation_1
    form.organisation_2.data = info.organisation_2
    form.organisation_3.data = info.organisation_3
    return render_template("main/profile.html.j2", form=form)


@bp_main.route("/submitter")
@login_required
@auth_role("submitter")
def submitter():
    return render_template("main/submitter.html.j2", name=current_user.info.name)


@bp_main.route("/reviewer")
@login_required
@auth_role("reviewer")
def reviewer():
    return render_template("main/reviewer.html.j2", name=current_user.info.name)
