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

        response, status = Entry.submit(form.data)

        if status == 400:  # bad request
            flash("Could not process the request.", "error")
            return render_template("new/new_submit.html", form=form)
        if status == 409:  # conflict. entry already exists
            existing_accession = response["accession"]
            entry_url = url_for("edit.view_json", bgc_id=existing_accession)
            flash(
                f"An entry with a locus accession and coordinates that overlap already exists ({existing_accession})",
                "error",
            )
            return render_template("new/new_submit.html", form=form)

        as_task_id = response.get("status").get("id")

        flash("New entry submitted successfully.", "success")

        return redirect(url_for("antismash.as_status", as_task_id=as_task_id))

    return render_template("new/new_submit.html", form=form)


@bp_new.route("/new_mutation/<bgc_id>", methods=["GET", "POST"])
@login_required
def create_bgc_mutation(bgc_id: str):

    if request.method == "POST":
        response = Entry.mutate(bgc_id)

        if response.status_code != 200:
            flash(f"Error creating new mutation: {response.json()['error']}", "error")
            return redirect(url_for("main.main"))

        mutation_accession = response.json()["accession"]

        return redirect(url_for("edit.edit_bgc_redirect", bgc_id=mutation_accession))

    return render_template(
        "edit/new_mutation.html",
        bgc_id=bgc_id,
    )
