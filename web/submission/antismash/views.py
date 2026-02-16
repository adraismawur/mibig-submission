import datetime
import enum
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

ANTISMASH_STATE = {
	0: "Pending",
	1: "Downloading",
	2: "Running",
	3: "Failed",
	4: "Finished",
}


@bp_as.route("/status/<as_task_id>")
@login_required
def as_status(as_task_id: str) -> Union[str, response.Response]:
    """Page to show status of antiSMASH processing task on a new entry

    Args:
        as_task_id (str): Asynchronous task identifier

    Returns:
        str | Response: rendered template or redirect to edit_bgc overview
    """

    response = requests.get(
        f"{current_app.config['API_BASE']}/antismash?guid={as_task_id}",
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    response_data = response.json()

    readable_state = "Unknown"

    if "state" in response_data:
        
        readable_state = ANTISMASH_STATE[response_data["state"]]
        submitted_at = datetime.datetime.fromisoformat(response_data["submitted_at"])
        readable_date = submitted_at.strftime("%Y-%m-%d %H:%M:%S")


        time_span = datetime.datetime.now(submitted_at.tzinfo) - submitted_at
        # omit ms
        time_elapsed = str(time_span).split(".")[0]

        if response_data.get("state") == 4:
            accession = response_data.get("accession")
            bgc_id = response_data.get("bgc_id")
            # redirect using 303 to change POST to GET
            return redirect(
                url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="locitax"), code=303
            )

    return render_template(
        "antismash/status.html",
        as_task_id=as_task_id,
        status=response.json(),
        readable_state=readable_state,
        readable_date=readable_date,
        time_elapsed=time_elapsed,
    )


@bp_as.route("/view/<accession>", methods=["GET"])
@login_required
def as_view(accession: str) -> Union[str, response.Response]:
    """Redirect to antiSMASH results page for a given accession

    Args:
        accession (str): GenBank accession
    """

    return render_template("antismash/view.html", accession=accession)
