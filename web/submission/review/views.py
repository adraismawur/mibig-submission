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
    PENDING = "pending_review"
    ACCEPTED = "accepted"


@bp_review.route("/", methods=["GET", "POST"])
@login_required
def list_submissions():
    # get the list of submissions that are marked ready for review
    submissions_api_path = (
        f"{current_app.config['API_BASE']}/submission?state={SUBMISSION_STATE.PENDING}"
    )
    submissions = requests.get(submissions_api_path).json()

    return render_template("review/list_submissions.html", submissions=submissions)
