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


readable_category_map = {
    "locitax": "Loci and taxonomy information",
    "biosynth": "Biosynthetic information",
    "compounds": "Compound information",
    "gene_information": "Gene information",
    "finalize": "Completeness and embargo",
    "full": "Full entry",
}


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
    search = request.args.get("search") or ""
    category = request.args.get("category") or ""
    start = request.args.get("start") or 0
    limit = request.args.get("limit") or 10

    reviewing = requests.get(
        f"{current_app.config['API_BASE']}/reviews/active",
        headers={"Authorization": f"Bearer {session['token']}"},
    ).json()
    pending_response = requests.get(
        f"{current_app.config['API_BASE']}/reviews/pending?start={start}&limit={limit}&search={search}&category={category}",
        headers={"Authorization": f"Bearer {session['token']}"},
    ).json()

    pending_count = pending_response["review_count"]
    pending_submissions = pending_response["reviews"]

    return render_template(
        "review/list_submissions.html",
        pending_submissions=pending_submissions,
        pending_count=pending_count,
        reviewing=reviewing,
        start=start,
        limit=limit,
        search=search,
        category=category,
    )

@bp_review.route("/claim_review/<bgc_id>/<category>", methods=["GET", "POST"])
@login_required
def claim_review(bgc_id: str, category: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/claim_review/",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category
            }
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/claim_review.html", bgc_id=bgc_id, category=readable_category_map[category])

@bp_review.route("/cancel/<bgc_id>/<category>", methods=["GET", "POST"])
@login_required
def cancel_review(bgc_id: str, category: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/cancel_review/",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category
            }
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/cancel_review.html", bgc_id=bgc_id, category=readable_category_map[category])

@bp_review.route("/approve/<bgc_id>/<category>", methods=["GET", "POST"])
@login_required
def approve(bgc_id: str, category: str):
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/accept/",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category
            }
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template(
        "review/approve.html",
        bgc_id=bgc_id,
        readable_category=readable_category_map[category],
    )
