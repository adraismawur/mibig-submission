from flask import blueprints

bp_as = blueprints.Blueprint("antismash", __file__, url_prefix="/antismash")

from submission.antismash import views
