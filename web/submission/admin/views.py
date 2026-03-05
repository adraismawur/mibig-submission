from flask import (
    current_app,
    flash,
    request,
    render_template,
    redirect,
    session,
    url_for,
)
import requests
from sqlalchemy import or_


from submission.admin import bp_admin
from submission.auth import auth_role
from submission.models import Role, User, UserInfo

from submission.admin.forms import UserAdd, UserEdit
from submission.models.users import UserRole

USER_ENDPOINT = "/user"


@bp_admin.route("/")
def index() -> str:
    return render_template("admin/index.html.j2")


@bp_admin.route("/users", methods=["GET"])
def list_users() -> str:
    response = requests.get(
        f"{current_app.config['API_BASE']}" + USER_ENDPOINT,
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    users = [User.from_json(user) for user in response.json()]

    return render_template("admin/users.html.j2", users=users)


@bp_admin.route("/users", methods=["POST"])
def search_users() -> str:
    search = request.form["search"]

    response = requests.get(
        f"{current_app.config['API_BASE']}{USER_ENDPOINT}?search={search}",
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    users = [User.from_json(user) for user in response.json()]

    return render_template("admin/user_search.html.j2", users=users)


@bp_admin.route("/user/<user_id>", methods=["GET"])
def user(user_id: int) -> str:
    user = User.query.get_or_404(user_id)
    return render_template("admin/user_list_line.html.j2", user=user)


@bp_admin.route("/user/<user_id>/edit", methods=["GET", "PUT"])
def user_edit(user_id: int) -> str:
    user = User.get_user(user_id)
    all_roles = [role for role in Role]
    form = UserEdit(
        request.form,
        data={
            "email": user.email,
        },
    )
    form.roles.choices = [(role.value, role.value) for role in all_roles]

    if form.validate_on_submit():
        user.email = form.email.data
        user.active = form.active.data
        roles = []
        for wanted_role in form.roles.data:
            for role in all_roles:
                if role.value == wanted_role:
                    roles.append(UserRole.from_enum(role))
        user.roles = roles

        response = requests.patch(
            f"{current_app.config['API_BASE']}/user/{user.id}",
            headers={
                "Authorization": f"Bearer {session['token']}",
                "Content-Type": "application/json",
            },
            data=user.to_json(),
        )

        if response.status_code != 200:
            error = response.json()["error"]
            flash(f"Could not add user: {error}", "error")
            return redirect(url_for("admin.list_users"), code=302)

        return render_template("admin/user_list_line.html.j2", user=user)

    form.roles.data = [role for role in user.roles]
    form.active.data = user.active

    return render_template("admin/user_edit_form.html.j2", form=form, user=user)


@bp_admin.route("/user/new", methods=["GET", "POST"])
def user_create() -> str:
    form = UserAdd(request.form)
    all_roles = [role for role in Role]
    form.roles.choices = [(role.value, role.value) for role in all_roles]

    if form.validate_on_submit():
        name = form.name.data
        info = UserInfo(
            alias=UserInfo.generate_alias(),
            name=name,
            call_name=UserInfo.guess_call_name(name),
            orc_id="",
            organisation_1=form.affiliation.data,
            organisation_2="",
            organisation_3="",
            public=False,
        )

        roles = []

        for wanted_role in form.roles.data:
            for role in all_roles:
                if role.value == wanted_role:
                    roles.append(UserRole.from_enum(role))

        user = User(
            id=None,
            email=form.email.data,
            active=form.active.data,
            roles=roles,
            info=info,
        )

        # save user
        response = requests.put(
            f"{current_app.config['API_BASE']}/user",
            headers={
                "Authorization": f"Bearer {session['token']}",
                "Content-Type": "application/json",
            },
            data=user.to_json(),
        )
        if response.status_code != 200:
            flash(f"Error creating user: {response.json()['error']}", "error")
            return render_template("admin/user_add.html.j2", form=form)

        return redirect(url_for("admin.list_users"), code=302)

    return render_template("admin/user_add.html.j2", form=form)
