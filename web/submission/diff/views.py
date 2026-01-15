from typing import Union
from flask import current_app, render_template, request
from flask_login import login_required
from werkzeug.wrappers import response

from submission.diff import bp_diff


@bp_diff.route("/", methods=["GET"])
@login_required
def view_diff() -> Union[str, response.Response]:
    left = request.args.get("left")
    right = request.args.get("right")

    mibig_json_url = current_app.config["API_BASE"] + "/entry/"

    return render_template(
        "diff/diff.html", left=left, right=right, mibig_json_url=mibig_json_url
    )
