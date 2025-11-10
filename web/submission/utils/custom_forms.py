"""Collection of custom form classes used throughout the submission system"""

import keyword

from markupsafe import Markup
from wtforms import (
    FieldList,
    Form,
    FormField,
    IntegerField,
    SelectField,
    StringField,
    validators,
)

from .custom_fields import ReferenceField
from .custom_widgets import FieldListAddBtn, TextInputWithSuggestions, SelectDefault
from .custom_validators import ValidateCitations, validate_loci


def location_form_factory(required: bool = False):
    """Create customized location form

    Args:
        required (bool): flag to add the InputRequired validator

    Returns:
        LocationForm: customized location form
    """

    if required:
        valids = [validators.InputRequired(), validate_loci]
    else:
        valids = [validators.Optional(), validate_loci]

    class LocationForm(Form):
        """Subform for location entry, use in combination with FormField"""

        pass

    setattr(LocationForm, "from", IntegerField("From", validators=valids))
    setattr(LocationForm, "to", IntegerField("To", validators=valids))

    return LocationForm


class ReferenceForm(Form):
    references = FieldList(StringField())


class EvidenceForm(Form):
    method = SelectField(
        "Method *",
        choices=[
            "In vitro expression",
            "Heterologous expression",
            "Enzymatic assays",
            "Knock-out studies",
            "Gene expression correlated with compound production",
            "Correlation of genomic and metabolomic data",
            "Synthetic-bioinformatic natural product (syn-BNP)",
            "Homology-based prediction",
        ],
        widget=SelectDefault(),
        validators=[validators.InputRequired()],
    )
    references = FieldList(
        StringField(),
        label="Citation(s) *",
        separator=" ",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        widget=FieldListAddBtn(
            label="Add reference",
        ),
    )


class SubtrateEvidenceForm(Form):
    name = SelectField(
        "Method *",
        choices=[
            "Activity assay",
            "ACVS assay",
            "ATP-PPi exchange assay",
            "Enzyme-coupled assay",
            "Feeding study",
            "Heterologous expression",
            "Homology",
            "HPLC",
            "In-vitro experiments",
            "Knock-out studies",
            "Mass spectrometry",
            "NMR",
            "Radio labelling",
            "Sequence-based prediction",
            "Steady-state kinetics",
            "Structure-based inference",
            "X-ray crystallography",
        ],
        widget=SelectDefault(),
        validate_choice=False,
        validators=[validators.InputRequired()],
    )
    references = ReferenceField(
        "Citation(s) *",
        widget=TextInputWithSuggestions(post_url="/edit/get_db_references"),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class StructureEvidenceForm(Form):
    method = SelectField(
        "Method *",
        choices=[
            "NMR",
            "Mass spectrometry",
            "MS/MS",
            "X-ray crystallography",
            "Chemical derivatisation",
            "Total synthesis",
            "MicroED",
            "Experimental values match with authentic standard",
        ],
        description="Technique used to elucidate/verify the structure",
        widget=SelectDefault(),
        validate_choice=False,
        validators=[validators.InputRequired()],
    )
    references = ReferenceField(
        "Citation(s) *",
        description="Comma separated list of references on this compound using this method",
        widget=TextInputWithSuggestions(post_url="/edit/get_db_references"),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class FunctionEvidenceForm(Form):
    method = SelectField(
        "Method *",
        choices=[
            "Other in vivo study",
            "Heterologous expression",
            "Knock-out",
            "Activity assay",
        ],
        widget=SelectDefault(),
        validate_choice=False,
        validators=[validators.InputRequired()],
    )
    references = ReferenceField(
        "Citation(s) *",
        widget=TextInputWithSuggestions(post_url="/edit/get_db_references"),
        validators=[validators.InputRequired(), ValidateCitations()],
    )
