from flask import current_app, session
from markupsafe import Markup
import requests
from wtforms import (
    DecimalField,
    Form,
    HiddenField,
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
from wtforms.widgets import HiddenInput
from submission.utils.custom_fields import ReferenceField, TagListField
from submission.utils.custom_forms import location_form_factory, LociEvidenceForm
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
    db_id = IntegerField(widget=HiddenInput())
    compound_id = IntegerField(widget=HiddenInput())
    name = StringField()
    observed = BooleanField()
    references = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class CompoundEvidence(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0, validators=None)
    db_compound_id = IntegerField(widget=HiddenInput(), default=0)
    method = StringField()
    references = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class CompoundsSubForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    name = StringField()
    evidence = FieldList(
        FormField(CompoundEvidence),
        widget=FieldListAddBtn(
            label="Add evidence",
        ),
        min_entries=1
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
