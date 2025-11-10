from typing import Any, Union

from flask import current_app, session
import requests
from sqlalchemy import select
from sqlalchemy.orm import Mapped, mapped_column, relationship

from submission.extensions import db
from submission.models import Reference


# TODO: expand table
class Entry(db.Model):
    __tablename__ = "entries"
    __table_args__ = {"schema": "edit"}

    id: Mapped[int] = mapped_column(autoincrement=True, primary_key=True)
    identifier: Mapped[str]
    references: Mapped[list["Reference"]] = relationship(
        "Reference",
        secondary="edit.entry_references",
        back_populates="entries",
        lazy="selectin",
    )

    @classmethod
    def create(cls, bgc_id: str):
        """Create a database entry for this BGC

        Args:
            bgc_id (str): BGC identifier
        """
        entry = cls(identifier=bgc_id)
        db.session.add(entry)
        db.session.commit()
        return entry

    def add_references(self, refs: list["Reference"]):
        """Add references to this entry if they are not already present

        Args:
            refs (list["Reference"]): list of Reference database objects
        """
        for ref in refs:
            if ref not in self.references:
                self.references.append(ref)
        db.session.commit()

    @staticmethod
    def get(bgc_id: str) -> Union["Entry", None]:
        """Get an entry from database based on identifier

        Args:
            bgc_id (str): BGC identifier

        Returns:
            Entry | None: entry database object or none if not exists
        """
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            entry = response.json()

            # we have a python keyword conflict and no good way to tell the form to
            # behave, so we have to mess with the json
            for i in range(len(entry["loci"])):
                entry["loci"][i]["location"]["from_"] = entry["loci"][i]["location"][
                    "from"
                ]

            # and another
            if entry["biosynthesis"]["classes"] is not None:
                for i in range(len(entry["biosynthesis"]["classes"])):
                    c = entry["biosynthesis"]["classes"][i]["class"]
                    entry["biosynthesis"]["classes"][i]["class_"] = c

            return entry
        return None
    
    def get_module(bgc_id: str, name: str):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module/{name}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            entry = response.json()

            return entry
        return None

    @staticmethod
    def get_or_create(bgc_id: str) -> "Entry":
        """Get an entry from the database or create one if it does not exist

        Args:
            bgc_id (str): BGC identifier

        Returns:
            Entry: new or existing database entry
        """
        if bgc_id == "new":
            entry = Entry.create(bgc_id=bgc_id)
        # entry = Entry.get(bgc_id=bgc_id)
        return entry

    # TODO: save all important data
    @staticmethod
    def submit(data: dict[str, Any]) -> str:
        """Submit a new entry to the API

        Args:
            bgc_id (str): BGC identfier
            data (dict): Minimal information to save
        """

        # we need to fix the "from_" attribute to lose the underscore
        # this is because WTForms does not have a good way (that I can
        # find) of naming a field something that is also a keyword in
        # python. Sigh

        for i in range(len(data["loci"])):
            data["loci"][i]["location"]["from"] = data["loci"][i]["location"]["from_"]
            del data["loci"][i]["location"]["from_"]

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )

        if response.status_code != 200:
            return None

        return response.json()

    @staticmethod
    def save_structure(bgc_id: str, data: dict[str, Any]):
        """Save structure information

        Args:
            bgc_id (str): BGC identfier
            data (dict): structure information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for structure in data["structures"]:
            for evidence in structure["evidence"]:
                refs.update(evidence["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_activity(bgc_id: str, data: dict[str, Any]):
        """Save activity information

        Args:
            bgc_id (str): BGC identfier
            data (dict): activity information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for activity in data["activities"]:
            if assays := activity.get("assays"):
                for assay in assays:
                    refs.update(assay["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_biosynth(bgc_id: str, b_class: str, data: dict[str, Any]):
        """Save biosynth information

        Args:
            bgc_id (str): BGC identfier
            b_class (str): biosynthetic class
            data (dict): biosynth information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if b_class == "NRPS":
            for rel_type in data["release_types"]:
                if references := rel_type.get("references"):
                    refs.update(references)
        elif b_class == "Saccharide":
            for glyc in data["glycosyltransferases"]:
                refs.update(glyc["references"])
            if subcl := data.get("subclusters"):
                for sub in subcl:
                    if references := sub.get("references"):
                        refs.update(references)

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_biosynth_paths(bgc_id: str, data: dict[str, Any]):
        """Save biosynthetic path information

        Args:
            bgc_id (str): BGC identfier
            data (dict): biosynthetic path information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for path in data["paths"]:
            refs.update(path["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_tailoring(bgc_id: str, data: dict[str, Any]):
        """Save tailoring information

        Args:
            bgc_id (str): BGC identfier
            data (dict): tailoring information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if enzymes := data.get("enzymes"):
            for enzyme in enzymes:
                refs.update(enzyme["enzyme"][0]["references"])
                if reactions := enzyme.get("reactions"):
                    for reaction in reactions:
                        for evidence in reaction["reaction_smarts"][0]["evidence_sm"]:
                            refs.update(evidence["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_annotation(bgc_id: str, data: dict[str, Any]):
        """Save annotation information

        Args:
            bgc_id (str): BGC identfier
            data (dict): annotation information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if annotations := data.get("annotations"):
            for annotation in annotations:
                if functions := annotation.get("functions"):
                    for function in functions:
                        if refers := function["mutation_phenotype"].get("references"):
                            refs.update(refers)
                        for evidence in function["evidence"]:
                            refs.update(evidence["references"])

        if domains := data.get("domains"):
            for domain in domains:
                if substrates := domain.get("substrates"):
                    for substrate in substrates:
                        if evidences := substrate.get("evidence"):
                            for evidence in evidences:
                                refs.update(evidences["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    # TODO: create batch entries from data dir
