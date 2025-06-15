# Task Tracking Service

We'll expose this as a set of tools for the agent to use. The actions will include:

* list items for conversation
* create item
* delete item

Items will have a:

* name: e.g. `Purchase Breakfast for 8/23 Campout"
* description: optional, e.g. `The campout on 8/23 will include a community breakfast. Joe will purchase the items at Costco.`
* source: e.g. `2025-06-30 Joe said "we need to get breakfast for everyone, I'll run by Costco on Friday and pick stuff up"`
* status: open; canceled; completed on {DATE}
* due date: optional, {DATE}
