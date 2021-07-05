import os
import pymongo
import json
import validators
from uuid import uuid4
from datetime import datetime
from dotenv import load_dotenv, find_dotenv
from flask import Flask, redirect, url_for, request, Response
from pymongo import MongoClient

app = Flask(__name__)


###### Database connection
load_dotenv(find_dotenv())
connection_string = os.getenv("CONNECTION_STRING")

cluster = MongoClient(connection_string)
db = cluster["url-shortener"]
collection = db["urls"]


###### Functions
def add_url(short_url, long_url):
    # Check if long_url is at the right format
    if long_url[:7] != "http://" and long_url[:8] != "https://":
        long_url = "http://" + long_url

    if not validators.url(long_url):
        response = Response(
            json.dumps(
                {
                    "status": 400,
                    "message": "Bad request: The provided long url's format is invalid.",
                }
            ),
            status=400,
            mimetype="application/json",
        )
        return response

    # If short_url is provided, check if it already exists
    if short_url != None:
        result = collection.find_one({"short_url": short_url})
        if result != None:
            response = Response(
                json.dumps(
                    {
                        "status": 400,
                        "message": "Bad request: The provided short url already exists.",
                    }
                ),
                status=400,
                mimetype="application/json",
            )
            return response
    # Else generate random short_url
    else:
        short_url = str(uuid4())[:4]

    # Add short url if it doesn't already exist
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


def get_url(short_url):
    # Return error if short url doesn't exist
    result = collection.find_one({"short_url": short_url})
    if result == None:
        return "Error: short_url doesn't exist"

    # Return the short url if it exists
    return result


###### Routes
@app.route("/", methods=["GET"])
def default():
    response = Response(
        "{'status': '200', 'message': Api is live. Read the documentation at ###'}",
        status=200,
        mimetype="application/json",
    )
    response.headers["Access-Control-Allow-Origin"] = "*"
    return response


@app.route("/shorten", methods=["GET"])
def shorten():
    short_url = request.args.get("short")
    long_url = request.args.get("long")

    response = add_url(short_url, long_url)
    return response


@app.route("/retrieve", methods=["GET"])
def retrieve():
    return


if __name__ == "__main__":
    app.run(debug=True)