from flask import current_app, session
from markupsafe import Markup
import requests
from wtforms import (
    DecimalField,
    Form,
    StringField,
    IntegerField,
    FieldList,
    SelectField,
    SelectMultipleField,
    FormField,
    BooleanField,
    ValidationError,
    validators,
    SubmitField,
)
from submission.utils.custom_fields import ReferenceField, TagListField
from submission.utils.custom_forms import location_form_factory, EvidenceForm
from submission.utils.custom_widgets import (
    FieldListAddBtn,
    TextInputIndicator,
    SubmitIndicator,
    SelectDefault,
    ProductInputSearch,
)
from submission.utils.custom_validators import (
    ValidateCitations,
    ValidateSingleInput,
    validate_genbank,
    validate_loci,
)


class BioactivitySubForm(Form):
    name = StringField()
    observed = BooleanField()
    references = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add reference",
        ),
    )


class CompoundsSubForm(Form):
    name = StringField()
    evidence = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )
    bioactivities = FieldList(
        FormField(BioactivitySubForm),
        widget=FieldListAddBtn(
            label="Add bioactivity",
        ),
    )
    structure = StringField()
    databaseIds = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add database",
        ),
    )
    mass = DecimalField("Mass")
    formula = StringField("Formula")


class CompoundsForm(Form):
    compounds = FieldList(
        FormField(CompoundsSubForm),
        widget=FieldListAddBtn(
            label="Add compound",
        ),
    )
