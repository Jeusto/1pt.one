import os
import pymongo
import json
import validators
from uuid import uuid4
from datetime import datetime
from dotenv import load_dotenv, find_dotenv
from flask import Flask, redirect, url_for, request, Response
from flask import render_template
from pymongo import MongoClient

app = Flask(__name__)

######### Database connection
load_dotenv(find_dotenv())
connection_string = os.getenv("CONNECTION_STRING")

cluster = MongoClient(connection_string)
database = cluster["url-shortener"]
collection = database["urls"]

######### Functions
# Create new short url and add it to the database
def add_url(short_url, long_url):
    # Check if long_url is at the right format
    if long_url[:7] != "http://" and long_url[:8] != "https://":
        long_url = "http://" + long_url

    if not validators.url(long_url):
        response = Response(
            json.dumps(
                {
                    "status": 400,
                    "message": "The provided long url's format is invalid.",
                }
            ),
            status=400,
            mimetype="application/json",
        )
        return response
    # If short_url is provided, check if it already exists
    if short_url != None and short_url != "":
        result = collection.find_one({"short_url": short_url})
        if result != None:
            response = Response(
                json.dumps(
                    {
                        "status": 400,
                        "message": "The provided short url already exists.",
                    }
                ),
                status=400,
                mimetype="application/json",
            )
            return response
    # If not, generate random short_url
    else:
        short_url = str(uuid4())[:4]

    # Add short url to database if it doesn't already exist
    current_time = datetime.now().strftime("%d/%m/%Y %H:%M:%S")
    unique_id = collection.estimated_document_count() + 1
    obj = {
        "_id": unique_id,
        "short_url": short_url,
        "long_url": long_url,
        "created_at": current_time,
        "number_of_visits": 0,
    }
    collection.insert_one(obj)
    response = Response(
        json.dumps(
            {
                "status": 201,
                "message": "Succesfully added short url.",
                "short_url": short_url,
                "long_url": long_url,
            }
        ),
        status=201,
        mimetype="application/json",
    )
    return response


# Retrieve info about a short url from database
def get_url(short_url):
    # Check if short url already exists
    result = collection.find_one({"short_url": short_url})
    if result == None:
        response = Response(
            json.dumps(
                {"status": 404, "message": "The provided short url does not exist."}
            ),
            status=404,
            mimetype="application/json",
        )
        return response

    # Return the short url if it exists
    response = Response(
        json.dumps(
            {
                "status": 200,
                "_id": result["_id"],
                "short_url": result["short_url"],
                "long_url": result["long_url"],
                "created_at": result["created_at"],
                "number_of_visits": result["number_of_visits"],
            }
        ),
        status=200,
        mimetype="application/json",
    )
    # Increase the number of visits
    collection.update(
        {"short_url": short_url},
        {"$set": {"number_of_visits": result["number_of_visits"] + 1}},
    )
    return response


######### Routes
# Default page
@app.route("/")
def index():
    return render_template("index.html")


# Show api status
@app.route("/status", methods=["GET"])
def status():
    response = Response(
        "{'status': '200', 'message': Api is live. Read the documentation at ###'}",
        status=200,
        mimetype="application/json",
    )
    response.headers["Access-Control-Allow-Origin"] = "*"
    return response


# Redirect from short url to long url
@app.route("/<short_url>")
def redirect_to_long(short_url):
    response = json.loads(get_url(short_url).data.decode())

    if response["status"] == 200:
        return redirect(response["long_url"])
    else:
        return redirect("https://google.com/")

    response = Response(
        "{'status': '200', 'message': Api is live. Read the documentation at ###'}",
        status=200,
        mimetype="application/json",
    )
    response.headers["Access-Control-Allow-Origin"] = "*"
    return response


# Create new short url
@app.route("/shorten", methods=["GET"])
def shorten():
    short_url = request.args.get("short")
    long_url = request.args.get("long")

    response = add_url(short_url, long_url)
    response.headers["Access-Control-Allow-Origin"] = "*"
    return response


# Retrieve info about a short url
@app.route("/retrieve", methods=["GET"])
def retrieve():
    short_url = request.args.get("short")

    response = get_url(short_url)
    response.headers["Access-Control-Allow-Origin"] = "*"
    return response


# Start server
if __name__ == "__main__":
    app.run(debug=True)