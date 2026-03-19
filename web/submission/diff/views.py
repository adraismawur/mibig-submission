from typing import Union
from flask import current_app, render_template, request
from flask_login import login_required
import requests
from werkzeug.wrappers import response

from submission.diff import bp_diff


@bp_diff.route("/", methods=["GET"])
@login_required
def view_diff() -> Union[str, response.Response]:
    left = request.args.get("left")
    right = request.args.get("right")

    mibig_json_url = current_app.config["API_BASE"] + "/export/entry"

    response_left = requests.get(f"{mibig_json_url}/{left}?pretty=true")
    response_right = requests.get(f"{mibig_json_url}/{right}?pretty=true")

    content_left = response_left.text
    content_right = response_right.text

    return render_template(
        "diff/diff.html", content_left=content_left, content_right=content_right, mibig_json_url=mibig_json_url
    )
