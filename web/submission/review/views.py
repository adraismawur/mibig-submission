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
from submission.utils.custom_forms import ReviewCommentForm


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


@bp_review.route("/accepted", methods=["GET", "POST"])
@login_required
def list_accepted_submissions():
    # get the list of submissions that are marked ready for review
    search = request.args.get("search") or ""
    category = request.args.get("category") or ""
    start = request.args.get("start") or 0
    limit = request.args.get("limit") or 10

    pending_response = requests.get(
        f"{current_app.config['API_BASE']}/reviews/accepted?start={start}&limit={limit}&search={search}&category={category}",
        headers={"Authorization": f"Bearer {session['token']}"},
    ).json()

    accepted_count = pending_response["review_count"]
    accepted_submissions = pending_response["reviews"]

    return render_template(
        "review/list_accepted_submissions.html",
        accepted_submissions=accepted_submissions,
        accepted_count=accepted_count,
        start=start,
        limit=limit,
        search=search,
        category=category,
    )

@bp_review.route("/review_comments/<bgc_id>/<category>")
@login_required
def view_review_comments(bgc_id: str, category: str):
    redirect = request.args.get('redirect')


    comment_endpoint = f"/review/{bgc_id}/{category}/comments"
    response = requests.get(
        f"{current_app.config['API_BASE']}" + comment_endpoint,
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    if response.status_code != 200:
        flash("Could not retrieve review comments", 'error')
        comments = []
    else:
        comments = response.json()

    return render_template("review/view_review_comments.html", bgc_id=bgc_id, comments=comments, redirect=redirect)

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

@bp_review.route("/rfc/<bgc_id>/<category>", methods=["GET", "POST"])
@login_required
def rfc(bgc_id: str, category: str):

    if request.form:
        form = ReviewCommentForm(request.form)
    else:
        form = ReviewCommentForm()
        
    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/rfc/",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category,
                "comment": form.data['comment']
            }
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template("review/rfc.html", bgc_id=bgc_id, category=readable_category_map[category], form=form)

@bp_review.route("/approve/<bgc_id>/<category>", methods=["GET", "POST"])
@login_required
def approve(bgc_id: str, category: str):

    if request.form:
        form = ReviewCommentForm(request.form)
    else:
        form = ReviewCommentForm()

    if request.method == "POST":
        response = requests.post(
            f"{current_app.config['API_BASE']}/submission/accept/",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category,
                "comment": form.data['comment']
            }
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("review.list_submissions"))

    return render_template(
        "review/approve.html",
        bgc_id=bgc_id,
        readable_category=readable_category_map[category],
        form=form,
    )

@bp_review.route("/references/<bgc_id>", methods=["GET"])
def view_references(bgc_id: str):
    redirect = request.args.get('redirect')
    
    response = requests.get(
        f"{current_app.config['API_BASE']}/review/{bgc_id}/references/",
        headers={"Authorization": f"Bearer {session['token']}"}
    )

    references = {}

    if response.status_code != 200:
        flash('Error getting references: ' + response.json()['error'])
    else:
        references = response.json()

    return render_template("review/list_references.html", references=references, redirect=redirect)

