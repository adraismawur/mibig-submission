from flask import (
    abort,
    current_app,
    render_template,
    render_template_string,
    request,
    redirect,
    session,
    url_for,
    flash,
)
from flask_login import current_user, login_required
import requests

from submission.review import bp_review
from submission.edit.forms.form_collection import FormCollection
from submission.models import Entry


class SUBMISSION_STATE:
    DRAFT = "draft"
    EDIT = "edit"
    PENDING = "pending review"
    REVIEWING = "being reviewed"
    ACCEPTED = "accepted"


@bp_review.route("/", methods=["GET", "POST"])
@login_required
def list_submissions():
    # get the list of submissions that are marked ready for review
    pending_submissions = requests.get(f"{current_app.config['API_BASE']}/submission?state={SUBMISSION_STATE.PENDING}").json()
    reviewing = requests.get(
        f"{current_app.config['API_BASE']}/reviews",
        headers={"Authorization": f"Bearer {session['token']}"},
    ).json()

    return render_template("review/list_submissions.html", pending_submissions=pending_submissions, reviewing=reviewing)

@bp_review.route("/claim_review/<bgc_id>", methods=["GET", "POST"])
@login_required
def claim_review(bgc_id: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/claim_review/{bgc_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/claim_review.html", bgc_id=bgc_id)

@bp_review.route("/cancel/<bgc_id>", methods=["GET", "POST"])
@login_required
def cancel_review(bgc_id: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/cancel_review/{bgc_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/cancel_review.html", bgc_id=bgc_id)

@bp_review.route("/approve/<bgc_id>", methods=["GET", "POST"])
@login_required
def approve(bgc_id: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/accept/{bgc_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/approve.html", bgc_id=bgc_id)
