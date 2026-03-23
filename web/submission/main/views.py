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
from flask_login import login_required, current_user, login_user, logout_user
import requests

from submission.main import bp_main
from submission.main.forms import SelectExisting, UserDetailsEditForm
from submission.auth import auth_role
from submission.models.users import Role, User, UserInfo
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
            return redirect(url_for("new.new_entry"))

        # edit valid existing entry
        if request.form.get("edit") and form.validate():
            bgc_id = form.accession.data
            Storage.create_entry_if_not_exists(bgc_id)

            return redirect(url_for("edit.edit_bgc"))

    # get the list of submissions for this user

    submission_start = request.args.get("submission_start") or 0

    submission_limit = request.args.get("submission_limit") or 10

    submission_search = request.args.get("submission_search") or ""


    submissions_api_path = (
        f"{current_app.config['API_BASE']}/submission?start={submission_start}&limit={submission_limit}&search={submission_search}"
    )
    response = requests.get(
        submissions_api_path,
        headers={"Authorization": f"Bearer {session['token']}"},
    ).json()

    existing_submissions = response["submissions"]
    submission_count = response["record_count"]


    # all entries

    entry_start = request.args.get("entry_start") or 0

    entry_limit = request.args.get("entry_limit") or 10

    entry_search = request.args.get("entry_search") or ""

    existing_entries_api_path = f"{current_app.config['API_BASE']}/entry?start={entry_start}&limit={entry_limit}&search={entry_search}"

    response = requests.get(existing_entries_api_path).json()

    existing_entries = response["entries"]
    entry_count = response["record_count"]

    return render_template(
        "main/index.html",
        form=form,
        user=current_user,
        existing_submissions=existing_submissions,
        submission_start=submission_start,
        submission_limit=submission_limit,
        submission_search=submission_search,
        submission_count=submission_count,
        existing_entries=existing_entries,
        entry_start=entry_start,
        entry_limit=entry_limit,
        entry_search=entry_search,
        entry_count=entry_count,
    )


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
    if request.form:
        form = UserDetailsEditForm(request.form)
    else:
        response = requests.get(
            f"{current_app.config['API_BASE']}/user/{current_user.id}",
            headers={"Authorization": "Bearer " + session["token"]},
        )

        if response.status_code != 200:
            flash("Could not get user information. Try logging in again")
            logout_user()
            return redirect(url_for("auth.login"))

        form = UserDetailsEditForm(data=response.json())
        
    if request.method == "POST":
        # current_user.info["name"] = request.form["name"]
        # current_user.info["call_name"] = request.form["call_name"]
        # current_user.info["orcid"] = request.form.get("orc_id", None)
        # current_user.info["organization1"] = request.form["organization1"]
        # current_user.info["organization2"] = request.form.get("organization2", None)
        # current_user.info["organization3"] = request.form.get("organization3", None)

        response = requests.patch(
            f"{current_app.config['API_BASE']}/user/{current_user.id}",
            json=form.data,
            headers={"Authorization": "Bearer " + session["token"]},
        )

        if response.status_code != 200:
            flash("Error updating your user details: " + str(response.json()))
            return render_template("main/profile.html")

        login_user(User.get_user(current_user.id))

        flash("Updated your user details")

        return redirect(url_for("main.profile"))
    


    return render_template("main/profile.html", form=form)


@bp_main.route("/submitter")
@login_required
@auth_role(Role.SUBMITTER)
def submitter():
    return render_template("main/submitter.html.j2", name=current_user.info.name)


@bp_main.route("/reviewer")
@login_required
@auth_role(Role.REVIEWER)
def reviewer():
    return render_template("main/reviewer.html.j2", name=current_user.info.name)
