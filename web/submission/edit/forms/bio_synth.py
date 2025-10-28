from flask import current_app, session
import requests
from wtforms import (
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
from submission.utils.custom_fields import TagListField
from submission.utils.custom_forms import location_form_factory, EvidenceForm
from submission.utils.custom_widgets import (
    FieldListAddBtn,
    TextInputIndicator,
    SubmitIndicator,
    SelectDefault,
    ProductInputSearch,
)
from submission.utils.custom_validators import (
    ValidateSingleInput,
    validate_genbank,
    validate_loci,
)


class BiosynthClassForm(Form):
    class_ = StringField(id="class", name="class")
    subclass = StringField()
    cyclases = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add cyclase",
        ),
        min_entries=0,
        default=[],
    )


class BiosynthForm(Form):
    classes = FieldList(
        FormField(BiosynthClassForm),
        min_entries=1,
        description="List of classes in this entry",
        widget=FieldListAddBtn(
            label="Add class",
        ),
    )


class BioSynthForm(Form):

    biosynthesis = FormField(BiosynthForm)

    submit = SubmitField("Continue to ...", widget=SubmitIndicator())
