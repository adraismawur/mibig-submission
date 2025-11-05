from dataclasses import dataclass
from submission.edit.forms.bio_synth import BioSynthForm
from submission.edit.forms.compounds import CompoundsForm
from submission.edit.forms.finalize import FinalizeForm
from submission.edit.forms.gene_annotation import GeneAnnotationForm
from submission.edit.forms.loci_tax import LociTaxonomyForm


@dataclass
class WizardPage:
    id: str
    form: type
    description: str

    def create_form(self, request_form, entry):
        return self.form(request_form, data=entry)


wizard_pages = [
    WizardPage("locitax", LociTaxonomyForm, "basic information"),
    WizardPage("biosynth", BioSynthForm, "biosynthetic gene information"),
    WizardPage("compounds", CompoundsForm, "compound information"),
    WizardPage("gene_annotation", GeneAnnotationForm, "additional gene annotation"),
    WizardPage("finalize", FinalizeForm, "final details"),
]

wizard_page_index = {page.id: i for i, page in enumerate(wizard_pages)}


def get_wizard_page(page_id: str) -> WizardPage:
    if page_id not in wizard_page_index:
        return None

    form_idx = wizard_page_index[page_id]

    if form_idx < 0 or form_idx > len(wizard_pages) - 1:
        return None

    return wizard_pages[form_idx]


def get_next_page(page_id: str) -> type:
    next_idx = wizard_page_index[page_id] + 1

    if next_idx > len(wizard_pages) - 1:
        return None

    return wizard_pages[next_idx]


def get_prev_page(page_id: str) -> tuple[int, type]:
    prev_idx = wizard_page_index[page_id] - 1

    if prev_idx < 0:
        return None

    return wizard_pages[prev_idx]
