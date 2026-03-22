from wtforms import (
    Form,
    IntegerField,
    StringField,
    HiddenField,
    BooleanField,
    FormField,
    SelectField,
    FieldList,
    validators,
)
from wtforms.widgets import HiddenInput
from submission.utils.custom_fields import (
    ReferenceField,
    TagListField,
    GeneIdField,
    smiles_field_factory,
)
from submission.utils.custom_forms import location_form_factory, SubtrateEvidenceForm
from submission.utils.custom_widgets import (
    TextInputWithSuggestions,
    SelectDefault,
    FieldListAddBtn,
)
from submission.utils.custom_validators import ValidateCitations


class CondensationDomain(Form):
    # "required": ["type", "gene", "location"]
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)
    _type = HiddenField("condensation")
    gene = GeneIdField("Gene")
    location = FormField(location_form_factory(location_default=-1))
    subtype = SelectField(
        "Subtype",
        choices=[
            "Dual",
            "Starter",
            "LCL",
            "DCL",
            "Ester bond-forming",
            "Heterocyclization",
        ],
        widget=SelectDefault(),
        validate_choice=False,
    )
    references = ReferenceField(
        "Citation(s)",
        widget=TextInputWithSuggestions(post_url="/edit/get_db_references"),
        validators=[ValidateCitations()],
    )


class AdenylationDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)

    # "required": ["type", "gene", "location"],
    class SubstateForm(Form):
        # "required": ["name", "proteinogenic", "structure"]
        name = StringField(
            "Name",
            widget=TextInputWithSuggestions(
                post_url="/edit/query_substrates",
                trigger="input changed delay:500ms, search, load",
            ),
        )
        proteinogenic = BooleanField("proteinogenic?")
        structure = smiles_field_factory(label="Structure (SMILES)")

    _type = HiddenField("adenylation")
    gene = GeneIdField("Gene")
    location = FormField(location_form_factory(location_default=-1))
    inactive = BooleanField("Inactive?")
    evidence = FieldList(
        FormField(SubtrateEvidenceForm),
        widget=FieldListAddBtn(label="Add additional evidence"),
    )
    precursor_biosynthesis = TagListField("Gene(s) involved in precursor biosynthesis")
    substrates = FieldList(
        FormField(SubstateForm),
        widget=FieldListAddBtn(label="Add additional substrate"),
    )


class CarrierDomain(Form):
    # "required": ["type", "gene", "location"]
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)
    _type = HiddenField("carrier")
    subtype = SelectField(
        "Subtype", choices=["ACP", "PCP"], widget=SelectDefault(), validate_choice=False
    )
    gene = GeneIdField("Gene")
    location = FormField(location_form_factory(location_default=-1))
    inactive = BooleanField("Inactive?")
    beta_branching = BooleanField("Beta-branching?")
    evidence = FieldList(
        FormField(SubtrateEvidenceForm),
        widget=FieldListAddBtn(label="Add additional evidence"),
    )


class AminotransferaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="aminotransferase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")
    inactive = BooleanField("Inactive?")


class CyclaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="cyclase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")
    references = ReferenceField(
        "Citation(s)",
        widget=TextInputWithSuggestions(post_url="/edit/get_db_references"),
        validators=[ValidateCitations()],
    )


class DehydrataseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="dehydratase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


class EnoylreductaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="enoylreductase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


class EpimeraseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="epimerase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


class HydroxylaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="hydroxylase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


class KetoreductaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="ketoreductase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")
    inactive = BooleanField("Inactive?")
    stereochemistry = SelectField(
        "Stereochemistry",
        choices=["A1", "A2", "B1", "B2", "C1", "C2"],
        widget=SelectDefault(),
        validate_choice=False,
    )
    evidence = FieldList(
        FormField(SubtrateEvidenceForm),
        widget=FieldListAddBtn(label="Add additional evidence"),
    )


class MethyltransferaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="methyltransferase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")
    subtype = SelectField(
        "Subtype",
        choices=["C", "N", "O", "other"],
        widget=SelectDefault(),
        validate_choice=False,
    )
    details = StringField("Details")


class OtherDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    # "required": ["type", "subtype", "gene", "location"]
    domain_type = StringField(widget=HiddenInput(), default="other")
    subtype = StringField("Subtype *", validators=[validators.InputRequired()])
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


class OxidaseDomain(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    domain_type = StringField(widget=HiddenInput(), default="oxidase")
    gene = GeneIdField("Gene *", validators=[validators.InputRequired()])
    location = FormField(location_form_factory(required=True, location_default=-1), label="Location *")


modification_domain_types = [
    "other",
    "epimerase",
    "ketoreductase",
    "enoylreductase",
    "oxidase",
    "aminotransferase",
    "branching",
    "dehydratase",
    "ligase",
    "hydroxylase",
    "product_template",
    "thioreductase",
    "methyltransferase",
    "thioesterase",
]


class ModificationDomainForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)
    type = SelectField(choices=modification_domain_types)
    gene = StringField()
    location = FormField(location_form_factory(location_default=-1))


# class ModificationDomainForm(Form):
#     aminotransferase = FieldList(
#         FormField(AminotransferaseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     cyclase = FieldList(
#         FormField(CyclaseDomain), widget=FieldListAddBtn(label="Add additional domain")
#     )
#     dehydratase = FieldList(
#         FormField(DehydrataseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     enoylreductase = FieldList(
#         FormField(EnoylreductaseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     epimerase = FieldList(
#         FormField(EpimeraseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     hydroxylase = FieldList(
#         FormField(HydroxylaseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     ketoreductase = FieldList(
#         FormField(KetoreductaseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     methyltransferase = FieldList(
#         FormField(MethyltransferaseDomain),
#         widget=FieldListAddBtn(label="Add additional domain"),
#     )
#     oxidase = FieldList(
#         FormField(OxidaseDomain), widget=FieldListAddBtn(label="Add additional domain")
#     )
#     other = FieldList(
#         FormField(OtherDomain), widget=FieldListAddBtn(label="Add additional domain")
#     )


class MonomerForm(Form):
    # "required": ["evidence", "name", "structure"]
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)
    name = StringField("Name *", validators=[validators.InputRequired()])
    structure = smiles_field_factory(label="Structure (SMILES)", required=True)
    evidence = FieldList(
        FormField(SubtrateEvidenceForm),
        min_entries=1,
        widget=FieldListAddBtn(label="Add additional evidence"),
    )

class AcyltransferaseForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_module_id = IntegerField(widget=HiddenInput(), default=0)

    class SubstrateForm(Form):
        name = SelectField(
            "Name",
            choices=["malonyl-CoA", "methylmalonyl-CoA", "ethylmalonyl-CoA", "other"],
            widget=SelectDefault(),
            validate_choice=False,
        )
        structure = smiles_field_factory(
            label="Structure (SMILES)",
            description="Please provide a substrate structure if you've selected 'other'.",
        )
        details = StringField("Details (Optional)")

    _type = HiddenField("acyltransferase")
    gene = GeneIdField()
    location = FormField(location_form_factory(location_default=-1))
    subtype = SelectField(
        choices=["cis-AT", "trans-AT"], widget=SelectDefault(), validate_choice=False
    )
    inactive = BooleanField("Inactive?")
    substrates = FieldList(
        FormField(SubstrateForm),
        widget=FieldListAddBtn(label="Add additional substrate"),
    )
    evidence = FieldList(
        FormField(SubtrateEvidenceForm),
        widget=FieldListAddBtn(label="Add additional evidence"),
    )
