import requests

from typing import Union
from flask import (
    current_app,
    render_template,
    request,
    redirect,
    session,
    url_for,
)
from flask_login import login_required
from werkzeug.wrappers import response

from submission.antismash import bp_as


@bp_as.route("/view/<accession>", methods=["GET"])
@login_required
def as_view(accession: str) -> Union[str, response.Response]:
    """Redirect to antiSMASH results page for a given accession

    Args:
        accession (str): GenBank accession
    """

    return render_template("antismash/view.html", accession=accession)
