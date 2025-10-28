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
        f"{current_app.config['API_BASE']}/antismash/{as_task_id}",
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    if response.json().get("state") == 4:
        response_data = response.json()
        accession = response_data.get("accession")
        bgc_id = response_data.get("bgc_id")
        # redirect using 303 to change POST to GET
        return redirect(url_for("new.new_loci_tax", bgc_id=bgc_id), code=303)

    return render_template(
        "antismash/status.html", as_task_id=as_task_id, status=response.json()
    )


@bp_new.route("/<bgc_id>/locitax", methods=["GET", "POST"])
@login_required
def new_loci_tax(bgc_id: str):
    """First page of the wizard. Contains loci and taxonomy information"""
    entry = Entry.get(bgc_id=bgc_id)

    if not entry:
        flash("Entry not found.", "danger")
        return render_template("new/new_review.html", entry=None)

    form = FormCollection.locitax(request.form, data=entry)

    if request.method == "POST" and form.validate():
        # response = Entry.submit(form.data)

        if not response:
            flash("Error submitting minimal information.", "danger")
            return render_template("new/new_review.html", form=form, entry=entry)

        return redirect(url_for("new.new_biosynth", bgc_id=bgc_id), code=303)

    return render_template("new/new_review.html", form=form, entry=entry, bgc_id=bgc_id)


@bp_new.route("/<bgc_id>/biosynth", methods=["GET", "POST"])
@login_required
def new_biosynth(bgc_id: str):
    """Second page of the wizard. This contains biosynthetic class information of the
    entry"""
    entry = Entry.get(bgc_id=bgc_id)

    if not entry:
        flash("Entry not found.", "danger")
        return render_template("new/new_review.html", entry=None)

    form = FormCollection.biosynth(request.form, data=entry)

    if request.method == "POST" and form.validate():
        # response = Entry.submit(form.data)

        if not response:
            flash("Error submitting minimal information.", "danger")
            return render_template("new/new_review.html", form=form, entry=entry)

        return redirect(url_for("new.new_biosynth"))

    return render_template("new/new_review.html", form=form, entry=entry, bgc_id=bgc_id)
