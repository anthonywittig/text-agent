# Task Tracking Service

We'll expose this as a set of tools for the agent to use. The actions will include:

* list items for conversation
* create item
* delete item

The items will be stored in DynamoDB.

Items will have a:

* id: UUID
* conversation ID: e.g. `+15554443333;+16665554444;+17776665555`
* name: e.g. `Purchase Breakfast for 8/23 Campout`
* description: optional, e.g. `The campout on 8/23 will include a community breakfast. Joe will purchase the items at Costco.`
* source: e.g. `2025-06-30 Joe said "we need to get breakfast for everyone, I'll run by Costco on Friday and pick stuff up"`
* status: open; canceled; completed on {DATE}
