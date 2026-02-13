from flask import Blueprint

bp_review = Blueprint("review", __file__, url_prefix="/review")

from submission.review import views
