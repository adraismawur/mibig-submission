import re
import csv
from pathlib import Path
from typing import Union

from sqlalchemy import select, or_
from flask import (
    abort,
    render_template,
    render_template_string,
    request,
    redirect,
    url_for,
    flash,
)
from flask_login import current_user, login_required
from werkzeug.datastructures import MultiDict
from werkzeug.wrappers import response
from wtforms.widgets import html_params
from markupsafe import Markup

from submission.extensions import db
from submission.edit import bp_edit
from submission.edit.forms.form_collection import FormCollection
from submission.edit.forms.edit_select import EditSelectForm
from submission.models.users import Role
from submission.utils import Storage, draw_smiles_svg, draw_smarts_svg
from submission.utils.custom_validators import is_valid_bgc_id
from submission.utils.custom_errors import ReferenceNotFound
from submission.models import Entry, NPAtlas, Substrate


@bp_edit.route("/<bgc_id>", methods=["GET", "POST"])
@login_required
def edit_bgc(bgc_id: str) -> Union[str, response.Response]:
    """Overview page with navigation to forms for entry sections

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered overview template or redirect section form
    """

    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    form = EditSelectForm(request.form)
    if request.method == "POST":
        for option, chosen in form.data.items():
            if chosen:
                return redirect(url_for(f"edit.edit_{option}", bgc_id=bgc_id))
    return render_template("edit/edit.html", form=form, bgc_id=bgc_id)


@bp_edit.route("/<bgc_id>/minimal", methods=["GET", "POST"])
@login_required
def edit_minimal(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter minimal entry information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered form template or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    # try to fill data from existing entry
    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Minimal"))
        form = FormCollection.minimal(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.minimal(request.form)
        # do not prefill reviewed checkbox on failed post
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_minimal(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Minimal", request.form, current_user)
            flash("Submitted minimal entry!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")
    return render_template(
        "edit/min_entry.html",
        form=form,
        bgc_id=bgc_id,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/structure", methods=["GET", "POST"])
@login_required
def edit_structure(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter structure information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered form template or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Structure"))
        form = FormCollection.structure(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.structure(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_structure(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Structure", request.form, current_user)
            flash("Submitted structure information!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")

    # on GET query db for any products already present
    else:
        # if one (or more) is present prefill
        try:
            products = next(
                csv.reader(
                    [
                        MultiDict(Storage.read_data(bgc_id).get("Minimal")).get(
                            "products"
                        )
                    ],
                    skipinitialspace=True,
                )
            )
        except:
            products = []

        for product in products:
            if product not in [struct.data.get("name") for struct in form.structures]:
                form.structures.append_entry(data={"name": product})
    return render_template(
        "edit/structure.html",
        form=form,
        bgc_id=bgc_id,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/render_smiles", methods=["POST"])
@login_required
def render_smiles() -> Union[str, response.Response]:
    origin = request.headers["Hx-Trigger-Name"]
    smiles_string = request.form.get(origin)

    if smiles_string is None or not (smiles := smiles_string.strip()):
        return ""

    return draw_smiles_svg(smiles)


@bp_edit.route("/<bgc_id>/bioact", methods=["GET", "POST"])
@login_required
def edit_activity(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter biological activity information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered form templare or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Bio_activity"))
        form = FormCollection.bioact(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.bioact(request.form)
        reviewed = False
    if request.method == "POST" and form.validate():
        try:
            Entry.save_activity(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Bio_activity", request.form, current_user)
            flash("Submitted activity information!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")
    else:
        # prefill compounds
        try:
            products = next(
                csv.reader(
                    [
                        MultiDict(Storage.read_data(bgc_id).get("Minimal")).get(
                            "products"
                        )
                    ],
                    skipinitialspace=True,
                )
            )
        except:
            products = []

        for product in products:
            if product not in [act.data.get("compound") for act in form.activities]:
                form.activities.append_entry(data={"compound": product})
    return render_template(
        "edit/biological_activity.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth", methods=["GET", "POST"])
@login_required
def edit_biosynth(bgc_id: str) -> str:
    """Selection overview page for class-specific biosynthesis forms

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str: rendered template
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    return render_template("edit/biosynthesis.html", bgc_id=bgc_id)


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


@bp_edit.route("/<bgc_id>/biosynth/<b_class>", methods=["GET", "POST"])
@login_required
def edit_biosynth_class(bgc_id: str, b_class: str) -> Union[str, response.Response]:
    """Form to enter class-specific biosynthesis information

    Args:
        bgc_id (str): BGC identifier
        b_class (str): Biosynthetic class

    Returns:
        str | Response: rendered template or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(
            Storage.read_data(bgc_id).get(f"BioSynth_{b_class}")
        )
        form = getattr(FormCollection, b_class)(data)
        reviewed = data.get("reviewed")
    else:
        form = getattr(FormCollection, b_class)(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_biosynth(bgc_id=bgc_id, b_class=b_class, data=form.data)
            Storage.save_data(bgc_id, f"BioSynth_{b_class}", request.form, current_user)
            flash(f"Submitted {b_class} biosynthesis information!")
            return redirect(url_for("edit.edit_biosynth", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")

    return render_template(
        "edit/biosynth_class_specific.html",
        form=form,
        b_class=b_class,
        bgc_id=bgc_id,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


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


@bp_edit.route("/<bgc_id>/biosynth/paths", methods=["GET", "POST"])
@login_required
def edit_biosynth_paths(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter biosynthetic path information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Biosynth_paths"))
        form = FormCollection.paths(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.paths(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_biosynth_paths(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Biosynth_paths", request.form, current_user)
            flash("Submitted biosynthetic path information!")
            return redirect(url_for("edit.edit_biosynth", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")

    return render_template(
        "edit/biosynth_paths.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/biosynth/modules", methods=["GET", "POST"])
@login_required
def edit_biosynth_modules(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter biosynthetic module information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_biosynth overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Biosynth_modules"))
        form = FormCollection.modules(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.modules(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        # TODO: save to db
        Storage.save_data(bgc_id, "Biosynth_modules", request.form, current_user)
        flash("Submitted biosynthetic module information!")
        return redirect(url_for("edit.edit_biosynth", bgc_id=bgc_id))

    return render_template(
        "edit/biosynth_modules.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/<bgc_id>/tailoring", methods=["GET", "POST"])
@login_required
def edit_tailoring(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter tailoring enzyme information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Tailoring"))
        form = FormCollection.tailoring(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.tailoring(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_tailoring(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Tailoring", request.form, current_user)
            flash("Submitted tailoring information!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")

    return render_template(
        "edit/tailoring.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/render_smarts", methods=["POST"])
@login_required
def render_smarts() -> Union[str, response.Response]:
    origin = request.headers["Hx-Trigger-Name"]
    smarts_string = request.form.get(origin)

    if smarts_string is None or not (smarts := smarts_string.strip()):
        return ""

    return draw_smarts_svg(smarts)


@bp_edit.route("/<bgc_id>/annotation", methods=["GET", "POST"])
@login_required
def edit_annotation(bgc_id: str) -> Union[str, response.Response]:
    """Form to enter gene annotation information

    Args:
        bgc_id (str): BGC identifier

    Returns:
        str | Response: rendered template or redirect to edit_bgc overview
    """
    if not is_valid_bgc_id(bgc_id):
        return abort(403, "Invalid existing entry!")

    if not request.form:
        data: MultiDict = MultiDict(Storage.read_data(bgc_id).get("Annotation"))
        form = FormCollection.annotation(data)
        reviewed = data.get("reviewed")
    else:
        form = FormCollection.annotation(request.form)
        reviewed = False

    if request.method == "POST" and form.validate():
        try:
            Entry.save_annotation(bgc_id=bgc_id, data=form.data)
            Storage.save_data(bgc_id, "Annotation", request.form, current_user)
            flash("Submitted annotation information!")
            return redirect(url_for("edit.edit_bgc", bgc_id=bgc_id))
        except ReferenceNotFound as e:
            flash(str(e), "error")

    return render_template(
        "edit/annotation.html",
        bgc_id=bgc_id,
        form=form,
        is_reviewer=current_user.has_role(Role.REVIEWER),
        reviewed=reviewed,
    )


@bp_edit.route("/add_field", methods=["POST"])
@login_required
def add_field() -> str:
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

    if (entry := Entry.get(bgc_id)) is not None:
        for ref in entry.references:
            options += li(ref.identifier, ref.summarize())
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
