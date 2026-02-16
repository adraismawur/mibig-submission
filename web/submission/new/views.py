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

        return redirect(url_for("antismash.as_status", as_task_id=as_task_id))

    return render_template("new/new_submit.html", form=form)
