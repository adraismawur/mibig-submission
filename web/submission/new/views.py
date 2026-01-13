import requests
from typing import Union
from flask import (
    current_app,
    render_template,
    request,
    redirect,
    session,
    url_for,
    flash,
)
from flask_login import login_required, current_user
from werkzeug.wrappers import response

from submission.edit.forms.form_collection import FormCollection
from submission.models.entries import Entry
from submission.new import bp_new


@bp_new.route("/", methods=["GET", "POST"])
@login_required
def new_entry():
    form = FormCollection.new(request.form)

    if request.method == "POST" and form.validate():

        response = Entry.submit(form.data)
        as_task_id = response.get("status").get("id")

        flash("New entry submitted successfully.", "success")

        return redirect(url_for("new.as_status", as_task_id=as_task_id))

    return render_template("new/new_submit.html", form=form)


@bp_new.route("/pending/<as_task_id>", methods=["GET"])
@login_required
def as_status(as_task_id: str) -> Union[str, response.Response]:
    """Page to show status of antiSMASH processing task on a new entry

    Args:
        as_task_id (str): Asynchronous task identifier

    Returns:
        str | Response: rendered template or redirect to edit_bgc overview
    """

    response = requests.get(
        f"{current_app.config['API_BASE']}/antismash?guid={as_task_id}",
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    if response.json().get("state") == 4:
        response_data = response.json()
        accession = response_data.get("accession")
        bgc_id = response_data.get("bgc_id")
        # redirect using 303 to change POST to GET
        return redirect(
            url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="locitax"), code=303
        )

    return render_template(
        "antismash/status.html", as_task_id=as_task_id, status=response.json()
    )
