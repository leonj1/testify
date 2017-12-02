#!/usr/bin/python

import requests
import unittest
import json
import uuid

host = "http://localhost:8743"


def my_random_string(string_length=10):
    """Returns a random string of length string_length."""
    random = str(uuid.uuid4()) # Convert UUID format to a Python string.
    random = random.upper() # Make all characters uppercase.
    random = random.replace("-","") # Remove the UUID '-'.
    return random[0:string_length] # Return the random string.


class MyTest(unittest.TestCase):
    # @unittest.skip("testing skipping")
    def test_health_check(self):
        r = requests.get("{}/public/health".format(host))
        self.assertEqual(r.status_code, 200)

    # @unittest.skip("testing skipping")
    def test_add_confession(self):
        payload = {
           "name": "dev01",
           "entity_type": "kubernetes cluster",
           "last_update": {
             "date": "2017-02-01 03:04:05 PM",
             "by": "someone",
             "status": "pass"
           },
           "journal": {
             "2017-02-01 03:04:05 PM": {
               "by": "someone",
               "checks": [
                 {
                   "name": "checked mounts",
                   "date": "2017-02-01 03:04:05 PM",
                   "status": "pass"
                 }
               ]
             }
           }
         }

        r = requests.post("{}/confessions".format(host), json=payload)
        self.assertEqual(r.status_code, 201)

    @unittest.skip("testing skipping")
    def test_add_hardware_with_tags(self):
        randomHost = my_random_string()
        payload = {
            "host": randomHost,
            "tags": {
               "environment": "dev"
            }
        }
        r = requests.post("{}/hardware".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        r = requests.get("{}/tags/environment".format(host))
        self.assertEqual(r.status_code, 200)
        # TODO Find a way to ignore create_date
        # self.assertEqual(r.content, "")
        payload = {
            "rack": "1"
        }
        r = requests.post("{}/tags/{}".format(host, randomHost), json=payload)
        self.assertEqual(r.status_code, 201)
        r = requests.get("{}/hardware/{}".format(host, randomHost))
        self.assertEquals(r.status_code, 200)
        r = requests.get("{}/tags/environment".format(host))
        self.assertEquals(r.status_code, 200)
        # print 'Content: {}'.format(r.content)

    @unittest.skip("testing skipping")
    def test_add_hardware_with_tags_and_services(self):
        randomHost = my_random_string()
        payload = {
            "host": randomHost,
            "services": [
                {
                    "name": "dockerhub",
                    "short_name": "dockerhub",
                    "repo": "",
                    "version": "1.0",
                    "docker": False
                }
            ],
            "tags": {
               "environment": "dev"
            }
        }
        r = requests.post("{}/hardware".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        # TODO Find a way to ignore create_date
        # self.assertEqual(r.content, "")
        service = {
            "name": "Apache Solr",
            "short_name": "solr",
            "repo": "",
            "version": "5.3",
            "docker": False
        }
        r = requests.post("{}/services/{}".format(host, randomHost), json=service)
        self.assertEqual(r.status_code, 201)
        # update the version for a service
        service = {
            "short_name": "solr",
            "version": "6.0",
        }
        r = requests.put("{}/services/{}".format(host, randomHost), json=service)
        self.assertEqual(r.status_code, 200)


if __name__ == '__main__':
    unittest.main()

