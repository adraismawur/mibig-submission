from datetime import datetime
import re
import csv
from pathlib import Path
from typing import Union

import requests
from sqlalchemy import select, or_
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
from werkzeug.datastructures import MultiDict
from werkzeug.wrappers import response
from wtforms.widgets import html_params
from markupsafe import Markup

from submission.edit.forms.biosynthesis import get_class_form
from submission.edit.forms.biosynthesis_modules import get_module_form
from submission.edit.forms.biosynthesis_paths import PathForm
from submission.edit.forms.compounds import CompoundsSubForm
from submission.extensions import db
from submission.edit import bp_edit
from submission.edit.forms.form_collection import FormCollection
from submission.edit.forms.edit_select import EditSelectForm
from submission.edit.forms.wizard import (
    get_default_wizard_page,
    get_wizard_page,
    get_prev_page,
    get_next_page,
)
from submission.models.users import Role
from submission.utils import Storage, draw_smiles_svg, draw_smarts_svg
from submission.utils.custom_validators import is_valid_bgc_id
from submission.utils.custom_errors import ReferenceNotFound
from submission.models import Entry, NPAtlas, Substrate



@bp_edit.route("/<bgc_id>")
@login_required
def edit_bgc_redirect(bgc_id: str):
    lock_endpoint = "/lock/list/" + bgc_id
    response = requests.get(
        f"{current_app.config['API_BASE']}" + lock_endpoint,
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    entry_locks = response.json()

    lock_info = {}

    for lock in entry_locks:
        lock_info[lock["category"]] = lock

    lock_keys = [
        "locitax",
        "biosynth",
        "compounds",
        "gene_information",
    ]

    # show finalize details if this is a full lock
    # or if this is a new entry, which is hacky and not great
    # TODO: fix this being hacky and not great
    if "full" in lock_info or bgc_id[0:3] == "new":
        lock_keys.append("finalize")

    readable_category_map = {
        "locitax": "Loci and taxonomy information",
        "biosynth": "Biosynthetic information",
        "compounds": "Compound information",
        "gene_information": "Gene information",
        "finalize": "Completeness and embargo",
        "full": "Full entry",
    }

    return render_template(
        "edit/edit.html",
        bgc_id=bgc_id,
        lock_keys=lock_keys,
        lock_info=lock_info,
        readable_category_map=readable_category_map,
    )


def generate_wizard_page(bgc_id: str, form_id: str, show_nav: bool):
    # instantiate associated form
    wizard_page = get_wizard_page(form_id)

    # try to fill data from existing entry
    data = wizard_page.get_data(bgc_id)

    form = None
    if wizard_page.form:
        form = wizard_page.create_form(request.form, data)
    else:
        form = wizard_page.create_form(None, data)

    prev_form = get_prev_page(form_id)
    next_form = get_next_page(form_id)

    if request.method == "POST" and form.validate():
        try:
            success, response = wizard_page.post_data(bgc_id, form.data)
            if success:
                flash("Updated submission data")
            else:
                error = response["error"]
                flash(f"Failed to update submission data: {error}", "error")
        except ReferenceNotFound as e:
            flash(str(e), "error")

        if wizard_page.post_redirect:
            return redirect(url_for(**wizard_page.post_redirect))
        

        data = wizard_page.get_data(bgc_id)
        if wizard_page.form:
            form = wizard_page.create_form(request.form, data)
        else:
            form = wizard_page.create_form(None, data)

    # get list of antismash accessions associated with this entry
    antismash_list_endpoint = "/antismash/list/"
    response = requests.get(
        f"{current_app.config['API_BASE']}" + antismash_list_endpoint + bgc_id,
        headers={"Authorization": f"Bearer {session['token']}"},
    )
    if response.status_code == 200:
        antismash_accessions = response.json()

    antismash_json_url = current_app.config["API_BASE"] + "/antismash/json/"

    return render_template(
        wizard_page.template,
        form=form,
        bgc_id=bgc_id,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=False,
        show_nav=show_nav,
        next_form=next_form,
        prev_form=prev_form,
        form_description=wizard_page.description,
        antismash_json_url=antismash_json_url,
        antismash_accessions=antismash_accessions,
    )


@bp_edit.route("/<bgc_id>/<form_id>", methods=["GET", "POST"])
@login_required
def edit_bgc(bgc_id: str, form_id: str) -> Union[str, response.Response]:
    """Form to enter minimal entry information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered form template or redirect to edit_bgc overview
    """
    # lock checks
    lock_response = Entry.check_lock(bgc_id, form_id)

    if lock_response.status_code != 200:
        flash(f"Locking error: {lock_response.json()['error']}", "error")
        return redirect(url_for("edit.edit_bgc_redirect", bgc_id=bgc_id))

    show_nav = lock_response.json()["full"]

    return generate_wizard_page(bgc_id, form_id, show_nav)


@bp_edit.route("/<bgc_id>/lock/request/<category>", methods=["GET", "POST"])
@login_required
def request_lock(bgc_id: str, category: str):
    # lock checks
    if request.method == "POST":
        lock_response = Entry.request_lock(bgc_id, category)

        if lock_response.status_code != 200:
            flash(
                f"Could not request lock for category: {lock_response.json()['error']}",
                "error",
            )
            return redirect(url_for("edit.edit_bgc_redirect", bgc_id=bgc_id))
        else:
            flash(
                f"Lock requested successfully. Your lock will be removed at {lock_response.json()['unlocks_at']}"
            )
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id=category))

    return render_template("lock/lock_request.html", bgc_id=bgc_id, category=category)


@bp_edit.route("/<bgc_id>/lock/release/<category>", methods=["GET", "POST"])
@login_required
def release_lock(bgc_id: str, category: str):
    # lock checks
    if request.method == "POST":
        lock_response = Entry.release_lock(bgc_id, category)

        if lock_response.status_code != 200:
            flash(
                f"Could not release lock for category: {lock_response.json()['error']}",
                "error",
            )
        else:
            flash("Lock released successfully")

        return redirect(url_for("edit.edit_bgc_redirect", bgc_id=bgc_id))

    return render_template("lock/lock_release.html", bgc_id=bgc_id, category=category)


@bp_edit.route("/<bgc_id>/json", methods=["GET"])
@login_required
def view_json(bgc_id: str):
    json_endpoint = "/entry/" + bgc_id
    response = requests.get(
        f"{current_app.config['API_BASE']}" + json_endpoint + "?pretty=true",
        headers={"Authorization": f"Bearer {session['token']}"},
    )

    if response.status_code != 200:
        flash("Could not retrieve entry")

    entry = response.text

    return render_template("edit/view_json.html", entry=entry)


@bp_edit.route("/redraft/<bgc_id>", methods=["GET", "POST"])
@login_required
def redraft_bgc(bgc_id: str):

    entry_json = Entry.get_text(bgc_id)

    if request.method == "POST":
        submission_endpoint = "/submission/redraft/" + bgc_id
        response = requests.post(
            f"{current_app.config['API_BASE']}" + submission_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("main.index"))

    return render_template("edit/redraft.html", bgc_id=bgc_id, entry_json=entry_json)


@bp_edit.route("/promote/<bgc_id>", methods=["GET", "POST"])
@login_required
def promote_bgc(bgc_id: str):

    entry_json = Entry.get_text(bgc_id)

    if request.method == "POST":
        submission_endpoint = "/submission/promote/" + bgc_id
        response = requests.post(
            f"{current_app.config['API_BASE']}" + submission_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("main.index"))

    return render_template("edit/promote.html", bgc_id=bgc_id, entry_json=entry_json)


@bp_edit.route("/discard/<bgc_id>", methods=["GET", "POST"])
@login_required
def discard_bgc(bgc_id: str):

    entry_json = Entry.get_text(bgc_id)

    if request.method == "POST":
        submission_endpoint = "/submission/discard/" + bgc_id
        response = requests.post(
            f"{current_app.config['API_BASE']}" + submission_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            flash(response.json()["error"], "error")

        return redirect(url_for("main.index"))

    return render_template("edit/discard.html", bgc_id=bgc_id, entry_json=entry_json)


@bp_edit.route("/render_smiles_form", methods=["POST"])
@login_required
def render_smiles_form() -> Union[str, response.Response]:
    origin = request.headers["Hx-Trigger-Name"]
    smiles_string = request.form.get(origin)

    if smiles_string is None or not (smiles := smiles_string.strip()):
        return ""

    return draw_smiles_svg(smiles)


@bp_edit.route("/render_smiles", methods=["GET"])
@login_required
def render_smiles():
    smiles_string = request.args.get("smiles_string")
    return draw_smiles_svg(smiles_string)


@bp_edit.route("/class_buttons/<bgc_id>", methods=["POST"])
@login_required
def class_buttons(bgc_id: str) -> str:
    """Obtain buttons linking to relevant biosynthetic classes for this BGC

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str: html buttons linking to class-specific form templates
    """

    # grab classes for current bgc_id
    classes = MultiDict(Storage.read_data(bgc_id).get("Minimal")).getlist("b_class")

    if not classes:
        classes = ["NRPS", "PKS", "Ribosomal", "Saccharide", "Terpene", "Other"]
    class_btns = ""
    for cls in classes:
        class_btns += f"<a class='btn btn-light' style='margin: 5px' role='button' href='/edit/{bgc_id}/biosynth/{cls}'>{cls}</a>"
    return class_btns


@bp_edit.route("/<bgc_id>/biosynth/operons", methods=["GET", "POST"])
@login_required
def edit_biosynth_operons(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter operon information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Biosynth_operons"))
        form = FormCollection.operons(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.operons(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        # TODO: save to db
        Storage.save_data(bgc_id, "Biosynth_operons", request.form, current_user)
        flash("Submitted operon information!")
        return redirect(url_for("edit.edit_biosynth", bgc_id=bgc_id))

    return render_template(
        "edit/biosynth_operons.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth/new_class", methods=["GET"])
@login_required
def create_biosynth_class_new(bgc_id):
    return create_biosynth_class(bgc_id, None)


@bp_edit.route("/<bgc_id>/biosynth/new_class/<class_type>", methods=["GET", "POST"])
@login_required
def create_biosynth_class(
    bgc_id: str, class_type: str
) -> Union[str, response.Response]:
    """Create a new biosynthetic module for a certain BGC"""

    choices = [
        {"label": "NRPS", "value": "NRPS"},
        {"label": "PKS", "value": "PKS"},
        {"label": "Ribosomal", "value": "ribosomal"},
        {"label": "Saccharide", "value": "saccharide"},
        {"label": "Terpene", "value": "terpene"},
        {"label": "Other", "value": "class_other"},
    ]

    if class_type:
        form = get_class_form(class_type)(request.form)
    else:
        form = None

    if request.method == "POST" and form.validate():

        if Entry.create_class(bgc_id, form.data):
            flash("Created new biosynthetic module!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))
        else:
            flash("Error creating new biosynthetic module", "error")
            return redirect(
                url_for(
                    "edit.create_biosynth_class", bgc_id=bgc_id, class_type=class_type
                )
            )

    return render_template(
        "wizard/biosynth_class_new.html",
        bgc_id=bgc_id,
        form=form,
        choices=choices,
        class_type=class_type,
    )


@bp_edit.route(
    "/<bgc_id>/biosynth/edit_class/<class_id>/<class_type>", methods=["GET", "POST"]
)
@login_required
def edit_biosynth_class(
    bgc_id: str, class_id: str, class_type: str
) -> Union[str, response.Response]:
    """Form to enter biosynthetic class information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """

    if not request.form:
        classData = Entry.get_class(bgc_id, class_id)
        form = get_class_form(class_type)(data=classData)
        reviewed = False
    else:
        form = form = get_class_form(class_type)(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():

        form.data["class"] = form.data["class_"]

        error = Entry.update_class(bgc_id, class_id, form.data)
        if error:
            flash(f"Failed to update biosynthetic class: {error}", "error")
        else:
            flash("Updated biosynthetic class")

        return redirect(
            url_for(
                "edit.edit_biosynth_class",
                bgc_id=bgc_id,
                class_id=class_id,
                class_type=class_type,
            )
        )

    return render_template(
        "wizard/biosynth_class_edit.html",
        bgc_id=bgc_id,
        class_type=class_type,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth/delete_class/<class_id>", methods=["GET", "POST"])
@login_required
def remove_biosynth_class(bgc_id: str, class_id: int):
    # get pretty printed version of module data
    class_text = Entry.get_class_text(bgc_id, class_id)

    if request.method == "POST":
        Entry.delete_class(bgc_id, class_id)
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))

    return render_template(
        "wizard/biosynth_class_remove.html",
        bgc_id=bgc_id,
        class_id=class_id,
        class_text=class_text,
    )


@bp_edit.route("/<bgc_id>/biosynth/new_module", methods=["GET"])
@login_required
def create_biosynth_module_new(bgc_id):

    return create_biosynth_module(bgc_id, None)


@bp_edit.route("/<bgc_id>/biosynth/new_module/<module>", methods=["GET", "POST"])
@login_required
def create_biosynth_module(bgc_id: str, module: str) -> Union[str, response.Response]:
    """Create a new biosynthetic module for a certain BGC"""

    choices = [
        {"label": "Co-enzyme A ligase (CAL)", "value": "cal"},
        {"label": "NRPS Type I", "value": "nrps_type1"},
        {"label": "NRPS Type VI", "value": "nrps_type6"},
        {"label": "Iterative PKS", "value": "pks_iterative"},
        {"label": "Modular PKS", "value": "pks_modular"},
        {"label": "Trans-AT PKS", "value": "pks_trans_at"},
        {"label": "Other", "value": "module_other"},
    ]

    if module:
        form = get_module_form(module)(request.form)
    else:
        form = None

    if request.method == "POST" and form.validate():

        if Entry.create_module(bgc_id, form.data):
            flash("Created new biosynthetic module!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))
        else:
            flash("Error creating new biosynthetic module", "error")
            return redirect(
                url_for("edit.create_biosynth_module", bgc_id=bgc_id, module=module)
            )

    return render_template(
        "wizard/biosynth_module_new.html",
        bgc_id=bgc_id,
        form=form,
        choices=choices,
        module=module,
    )


@bp_edit.route(
    "/<bgc_id>/biosynth/edit_module/<module_id>/<module>", methods=["GET", "POST"]
)
@login_required
def edit_biosynth_module(
    bgc_id: str, module_id: int, module: str
) -> Union[str, response.Response]:
    """Form to enter biosynthetic module information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """

    if not request.form:
        moduleData = Entry.get_module(bgc_id, module_id)
        form = get_module_form(module)(data=moduleData)
        reviewed = False
    else:
        form = form = get_module_form(module)(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        if Entry.update_module(bgc_id, module_id, form.data):
            flash("Updated biosynthetic module!")
        else:
            flash("Failed to update biosynthetic module", "error")

        return redirect(
            url_for(
                "edit.edit_biosynth_module",
                bgc_id=bgc_id,
                module_id=module_id,
                module=module,
            )
        )

    return render_template(
        "wizard/biosynth_module_edit.html",
        bgc_id=bgc_id,
        module=module,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth/move_module/<id_from>/<id_to>")
@login_required
def move_biosynth_module(bgc_id: str, id_from: int, id_to: int):

    error = Entry.move_module(bgc_id, id_from, id_to)

    if error:
        flash("Could not reorder modules: " + str(error), "error")
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))

    return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))


@bp_edit.route("/<bgc_id>/biosynth/delete_module/<module_id>", methods=["GET", "POST"])
@login_required
def remove_biosynth_module(bgc_id: str, module_id: int):
    # get pretty printed version of module data
    module_text = Entry.get_module_text(bgc_id, module_id)

    if request.method == "POST":
        Entry.delete_module(bgc_id, module_id)
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))

    return render_template(
        "wizard/biosynth_module_remove.html",
        bgc_id=bgc_id,
        name=module_id,
        module_text=module_text,
    )


@bp_edit.route("/<bgc_id>/biosynth/new_path/<biosynth_id>/path", methods=["GET", "POST"])
@login_required
def create_biosynth_path(
    bgc_id: str, biosynth_id: int
) -> Union[str, response.Response]:
    """Create a new biosynthetic path for a certain BGC"""

    if request.form:
        form = PathForm(request.form)
    else:
        form = PathForm()

    if request.method == "POST" and form.validate():

        if Entry.create_path(bgc_id, int(biosynth_id), form.data):
            flash("Created new biosynthetic path!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))
        else:
            flash("Error creating new biosynthetic path", "error")
            return redirect(
                url_for(
                    "edit.create_biosynth_path", bgc_id=bgc_id, biosynth_id=biosynth_id
                )
            )

    return render_template(
        "wizard/biosynth_path_new.html",
        bgc_id=bgc_id,
        form=form,
    )


@bp_edit.route("/<bgc_id>/biosynth/edit_path/<path_id>/path", methods=["GET", "POST"])
@login_required
def edit_biosynth_path(bgc_id: str, path_id: int) -> Union[str, response.Response]:
    """Form to enter biosynthetic path information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """

    if not request.form:
        pathData = Entry.get_path(bgc_id, path_id)
        form = PathForm(data=pathData)
        reviewed = False
    else:
        form = PathForm(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        if Entry.update_path(bgc_id, path_id, form.data):
            flash("Updated biosynthetic path!")
        else:
            flash("Failed to update biosynthetic path", "error")

        return redirect(
            url_for(
                "edit.edit_biosynth_path",
                bgc_id=bgc_id,
                path_id=path_id,
            )
        )

    return render_template(
        "wizard/biosynth_path_edit.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth/delete_path/<path_id>", methods=["GET", "POST"])
@login_required
def remove_biosynth_path(bgc_id: str, path_id: int):
    # get pretty printed version of path data
    path_text = Entry.get_path_text(bgc_id, path_id)

    if request.method == "POST":
        Entry.delete_path(bgc_id, path_id)
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))

    return render_template(
        "wizard/biosynth_path_remove.html",
        bgc_id=bgc_id,
        path_id=path_id,
        path_text=path_text,
    )


@bp_edit.route("/<bgc_id>/compounds/new_compound", methods=["GET", "POST"])
@login_required
def create_compound(bgc_id: str):
    if request.form:
        form = CompoundsSubForm(request.form)
    else:
        form = CompoundsSubForm()

    if request.method == "POST":
        Entry.create_compound(bgc_id, form.data)
        flash("Compound created")
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="compounds"))

    return render_template(
        "wizard/compound_edit.html", bgc_id=bgc_id, form=form, new=True
    )


@bp_edit.route(
    "/<bgc_id>/compounds/<compound_id>/edit_compound", methods=["GET", "POST"]
)
@login_required
def edit_compound(bgc_id: str, compound_id: int):

    if not request.form:
        compoundData = Entry.get_compound(bgc_id, compound_id)

        form = CompoundsSubForm(data=compoundData["compounds"][0])
    else:
        form = CompoundsSubForm(request.form)

    if request.method == "POST":
        response = Entry.update_compound(bgc_id, form.data)
        if response.status_code == 200:
            flash("Compound updated successfully")
        else:
            flash("Failed to update compound: " + response.json()["error"], "error")

    return render_template(
        "wizard/compound_edit.html", bgc_id=bgc_id, form=form, new=False
    )


@bp_edit.route("/<bgc_id>/compounds/remove/<compound_id>", methods=["GET", "POST"])
@login_required
def remove_compound(bgc_id: str, compound_id: int):

    compound_text = Entry.get_compound_text(bgc_id, compound_id)

    if request.method == "POST":
        Entry.delete_compound(bgc_id, compound_id)
        return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="compounds"))

    return render_template(
        "wizard/compound_remove.html",
        bgc_id=bgc_id,
        compound_text=compound_text,
    )


def render_gene_information_edit(
    bgc_id: str, type: str, form_type: str, id: int, new: bool
):
    data_get_function = {
        "new_addition": Entry.get_gene_addition,
        "new_deletion": Entry.get_gene_deletion,
        "new_annotation": Entry.get_gene_annotation,
        "edit_addition": Entry.get_gene_addition,
        "edit_deletion": Entry.get_gene_deletion,
        "edit_annotation": Entry.get_gene_annotation,
    }

    data_set_function = {
        "new_addition": Entry.update_or_create_gene_addition,
        "new_deletion": Entry.update_or_create_gene_deletion,
        "new_annotation": Entry.update_or_create_gene_annotation,
        "edit_addition": Entry.update_or_create_gene_addition,
        "edit_deletion": Entry.update_or_create_gene_deletion,
        "edit_annotation": Entry.update_or_create_gene_annotation,
    }

    if new and request.method == "GET":
        form = getattr(FormCollection, form_type)()
    else:
        if request.method == "POST":
            form = getattr(FormCollection, form_type)(request.form)
            data, error = data_set_function[form_type](bgc_id, form.data)

            if error:
                flash(error, "error")
            else:
                flash("Update successful")

        else:
            data, error = data_get_function[form_type](bgc_id, id)

            if error:
                flash(error, "error")

            form = getattr(FormCollection, form_type)(data=data)

    return render_template(
        "wizard/gene_information_edit.html",
        form=form,
        new=new,
        bgc_id=bgc_id,
        type=type,
    )


@bp_edit.route("/<bgc_id>/save_operons", methods=["POST"])
@login_required
def save_operons(bgc_id: str):

    return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="biosynth"))


@bp_edit.route("/<bgc_id>/gene_information/new_addition", methods=["GET", "POST"])
@login_required
def add_gene_addition(bgc_id: str):
    return render_gene_information_edit(
        bgc_id, "gene addition", "new_addition", None, True
    )


@bp_edit.route("/<bgc_id>/gene_information/new_deletion", methods=["GET", "POST"])
@login_required
def add_gene_deletion(bgc_id: str):
    return render_gene_information_edit(
        bgc_id, "gene deletion", "new_deletion", None, True
    )


@bp_edit.route("/<bgc_id>/gene_information/new_annotation", methods=["GET", "POST"])
@login_required
def add_gene_annotation(bgc_id: str):
    return render_gene_information_edit(
        bgc_id, "gene annotation", "new_annotation", None, True
    )


@bp_edit.route(
    "/<bgc_id>/gene_information/<addition_id>/edit_addition", methods=["GET", "POST"]
)
@login_required
def edit_gene_addition(bgc_id: str, addition_id: int):
    return render_gene_information_edit(
        bgc_id, "gene addition", "edit_addition", addition_id, False
    )


@bp_edit.route(
    "/<bgc_id>/gene_information/<deletion_id>/edit_deletion", methods=["GET", "POST"]
)
@login_required
def edit_gene_deletion(bgc_id: str, deletion_id: int):
    return render_gene_information_edit(
        bgc_id, "gene addition", "edit_deletion", deletion_id, False
    )


@bp_edit.route(
    "/<bgc_id>/gene_information/<annotation_id>/edit_annotation",
    methods=["GET", "POST"],
)
@login_required
def edit_gene_annotation(bgc_id: str, annotation_id: int):
    return render_gene_information_edit(
        bgc_id, "gene addition", "edit_annotation", annotation_id, False
    )


@bp_edit.route(
    "/<bgc_id>/gene_information/delete/<type>/<information_id>", methods=["GET", "POST"]
)
@login_required
def remove_gene_information(bgc_id: str, type: str, information_id: int):

    types = {
        "gene_addition": {
            "human_readable": "gene addition",
            "get_method": Entry.get_gene_addition,
            "delete_method": Entry.remove_gene_addition,
        },
        "gene_deletion": {
            "human_readable": "gene deletion",
            "get_method": Entry.get_gene_deletion,
            "delete_method": Entry.remove_gene_deletion,
        },
        "gene_annotation": {
            "human_readable": "gene annotation",
            "get_method": Entry.get_gene_annotation,
            "delete_method": Entry.remove_gene_annotation,
        },
    }

    type_dict = types[type]

    if request.method == "POST":
        error = type_dict["delete_method"](bgc_id, information_id)

        if error is not None:
            flash(error, "error")

        return redirect(
            url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="gene_information")
        )

    data_text, error = type_dict["get_method"](bgc_id, information_id, pretty=True)

    if error is not None:
        flash(
            f"Could not retrieve data for {type_dict['human_readable']} to delete: {error}",
            "error",
        )
        return redirect(
            url_for("edit.edit_bgc", bgc_id=bgc_id, form_id="gene_information")
        )

    return render_template(
        "wizard/gene_information_remove.html",
        bgc_id=bgc_id,
        description=type_dict["human_readable"],
        data_text=data_text,
    )


@bp_edit.route("/render_smarts", methods=["POST"])
@login_required
def render_smarts() -> Union[str, response.Response]:
    origin = request.headers["Hx-Trigger-Name"]
    smarts_string = request.form.get(origin)

    if smarts_string is None or not (smarts := smarts_string.strip()):
        return ""

    return draw_smarts_svg(smarts)


@bp_edit.route("/add_field", methods=["POST"])
@login_required
def add_field(field=None) -> str:
    """Render an additional field as a subform

    Whenever a FieldList triggers this request for the addition of an entry to the list,
    use the id of the trigger to determine the where to append the entry.

    Examples:
        trigger id = 'structures'  -->   form.structures.append_entry()
        trigger id = 'enzymes-0-enzyme-0-auxiliary_enzymes'  -->
                        form.enzymes[0].enzyme[0].auxiliary_enzymes.append_entry()

    Returns:
        str: appended field rendered as a subform template string
    """

    if field:
        field.append_entry()
        return

    # find directions to field to append an entry to
    directions = request.headers["Hx-Trigger"].split("-")
    # get origin of request to determine which form to use
    formname = Path(request.referrer).name
    curr = getattr(FormCollection, formname)(request.form)

    # sequentially traverse form fields
    i = 0
    while i + 2 < len(directions):
        subform, subform_idx = directions[i : i + 2]

        # if fields were deleted, the index on the page will not match the actual index
        page_idx = set()
        base_field_name = "-".join(directions[: i + 1]) + "-"
        for key in request.form.keys():
            if key.startswith(base_field_name):
                page_idx.add(key.replace(base_field_name, "").split("-", 1)[0])
        if subform_idx in page_idx:
            actual_idx = sorted(page_idx).index(subform_idx)
        else:
            actual_idx = int(subform_idx)

        curr = getattr(curr, subform)[actual_idx]
        i += 2

    # until we reach the final field that issued the request
    final = getattr(curr, directions[i])
    final.append_entry()

    return render_template_string(
        """{% import 'macros.html' as m %}
        {{m.simple_divsubform(field, deletebtn=True)}}""",
        field=final[-1],
    )


@bp_edit.route("/get_db_references", methods=["POST"])
@login_required
def get_db_references() -> str:
    """Collect references connected to an entry and format into html list of suggestions

    Returns:
        str: HTML list of references
    """
    bgc_id_match = re.search("edit/([^/]+)/", request.referrer)
    if bgc_id_match is not None:
        bgc_id = bgc_id_match.group(1)

    li = (
        lambda val, descr: f"<li id={val} hx-post='/edit/append_reference' hx-target='previous input' hx-swap='outerHTML' hx-trigger='mousedown'>{descr}</li>"
    )
    options = (
        "<span class='text-muted form-text'>Known references for this entry:</span>"
    )

    # if (entry := Entry.get(bgc_id)) is not None:
    #     for ref in entry.references:
    #         options += li(ref.identifier, ref.summarize())
    return Markup(options)


@bp_edit.route("/append_reference", methods=["POST"])
@login_required
def append_reference() -> str:
    """Append a reference to an existing input

    Returns:
        str: a new html input element with an added reference
    """
    target = request.headers.get("Hx-Target")
    current = request.values.get(target)
    new_ref = request.headers.get("Hx-Trigger")

    if not current:
        current = f'"{new_ref}"'
    else:
        current_refs = set(next(csv.reader([current], skipinitialspace=True)))
        if new_ref not in current_refs:
            current += f', "{new_ref}"'
    return Markup(
        f"<input class='form-control' {html_params(value=current, id=target, name=target)}>"
    )


@bp_edit.route("/query_npatlas", methods=["POST"])
@login_required
def query_npatlas():
    """Check if a compound name is present in the NPAtas table and prefill information

    If no compound information is found, aborts the request.
    """
    trigger = request.headers.get("Hx-Trigger")
    compound = request.form.get(trigger)
    base = trigger.rpartition("-")[0]

    # if compound present in NPAtlas, prefill relevant data
    if (npa_entry := NPAtlas.get(compound)) is not None:
        # if entry exists, only overwrite npatlas data, keep the rest
        relevant_data = MultiDict(
            [(k, v) for k, v in request.form.items() if k.startswith(base)]
        )

        relevant_data.setlist(f"{base}-structure", [npa_entry.compound_smiles])
        relevant_data.setlist(f"{base}-formula", [npa_entry.compound_molecular_formula])
        relevant_data.setlist(f"{base}-mass", [str(npa_entry.compound_accurate_mass)])

        db_field = f"{base}-db_cross"
        npaid = f"npatlas:{npa_entry.npaid}"
        # taglistfield data is given as ['"data1", "data2", "data3"']
        if not (current_db_cross := relevant_data.getlist(db_field)[0]):
            relevant_data.setlist(db_field, [npaid])
        elif npaid not in current_db_cross:
            relevant_data.setlist(db_field, [f'{current_db_cross}, "{npaid}"'])

        form = FormCollection.structure(relevant_data)
        form.validate()
        return render_template_string(
            """{% import 'macros.html' as m %}
            {{m.simple_divsubform(field, deletebtn=true, message=message)}}""",
            field=form.structures[0],
            message=f"{compound} information filled from NPAtlas",
        )
    else:
        abort(404, "not present in npatlas")


@bp_edit.route("/render_npatlas_button", methods=["POST"])
@login_required
def render_npatlas_button() -> str:
    """Render a prefill button iff the entered compound is present in NPAtlas"""
    field_id = request.headers.get("Hx-Trigger")
    compound = request.form.get(field_id) or ""

    if NPAtlas.get(compound):
        render_kw = {
            "id": field_id,
            "name": field_id,
            "hx-post": "/edit/query_npatlas",
            "hx-target": "closest fieldset",
            "hx-swap": "outerHTML",
            "hx-trigger": "click",
        }
        add_btn = Markup(
            "<span class='fst-italic'>This compound is present in NPAtlas </span>"
            f"<button class='btn btn-light btn-sm form-text' {html_params(**render_kw)}>Fetch information</button>"
        )
        return add_btn
    return ""


@bp_edit.route("/query_product_name", methods=["POST"])
@login_required
def query_product():

    search_val = request.form.get("prod-search")
    if not search_val:
        return ""
    matching_compounds = db.session.scalars(
        select(NPAtlas).where(
            or_(
                NPAtlas.compound_names.istartswith(search_val),
                NPAtlas.npaid.istartswith(search_val),
            )
        )
    )
    result = [(res.compound_names, res.npaid) for res in matching_compounds]
    if not result:
        return Markup("<p class='form-text'><i>No matches found in NPAtlas</i></p>")

    row_kw = {
        "hx-post": "/edit/append_product",
        "hx-target": "previous #products",
        "hx-swap": "outerHTML",
    }

    row = (
        lambda name, idx: f"<tr {html_params(id=idx, **row_kw)}><td>{name}</td><td>{idx}</td></tr>"
    )
    return Markup(
        f"""<thead>
                  <tr>
                  <th>Compound Name</th>
                  <th>NPAtlas ID</th>
                  </tr>
                  </thead>
                  <tbody>
                  {"".join([row(n, i) for n,i in result])}
                  <tr></tr>
                  </tbody>
                  """
    )


@bp_edit.route("/append_product", methods=["POST"])
@login_required
def append_product():
    target = request.headers.get("Hx-Target")
    current = request.form.get(target)
    added_idx = request.headers.get("Hx-Trigger")
    entry = db.session.scalar(select(NPAtlas).where(NPAtlas.npaid == added_idx))

    if '"' in entry.compound_names:
        shown_name = entry.compound_names
    else:
        shown_name = f'"{entry.compound_names}"'

    if not current:
        new = shown_name
    else:
        new = current + f", {shown_name}"
    return Markup(
        f"<input class='form-control' {html_params(value=new, id=target, name=target)}>"
    )


@bp_edit.route("/query_substrates", methods=["POST"])
@login_required
def query_substrates() -> str:
    """Fetch a list of substrate options from the database

    Returns:
        str: collection of html list options
    """
    trigger = request.headers.get("Hx-Trigger")
    search_val = request.form.get(trigger)

    li = (
        lambda idx, descr: f"<li id={trigger} hx-post='/edit/fill_substrate/{idx}' hx-target='closest fieldset' hx-swap='outerHTML'>{descr}</li>"
    )

    options = "<span class='text-muted form-text'>Common substrates:</span>"
    for substrate in Substrate.isearch(search_val):
        options += li(substrate.id, substrate.summarize())
    return Markup(options)


@bp_edit.route("/fill_substrate/<idx>", methods=["POST"])
@login_required
def fill_substrate(idx: int) -> str:
    """Fill substrate information into form section

    Args:
        idx (str): db id of substrate entry

    Returns:
        str: rendered fieldset containing substrate information
    """
    trigger = request.headers.get("Hx-Trigger")
    base = trigger.rpartition("-")[0]

    if substrate := db.session.scalar(select(Substrate).where(Substrate.id == idx)):
        data = MultiDict(
            [
                (f"{base}-name", substrate.identifier),
                (f"{base}-structure", substrate.structure),
                (f"{base}-proteinogenic", substrate.proteinogenic),
            ]
        )
        form = FormCollection.modules(data)

        # both NPRS 1 and 6 have the same adenylation/substrate forms
        nrps_type = base.partition("-")[0]
        return render_template_string(
            """{% import 'macros.html' as m %}
                {{m.simple_divsubform(field, deletebtn=true)}}""",
            field=getattr(form, nrps_type)[0].a_domain[0].substrates[0],
        )
    else:
        abort(404, "substrate not found!")
