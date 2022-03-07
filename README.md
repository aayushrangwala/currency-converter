1. refresh every 5 minutes for exchange rates for all providers but with the first successful one only
- skip that provider if its not expired and use a new provider to add rates.
2. refresh on New() of the store 
3. refresh on every cache miss: code, provider. Considering the expiration time of each entry is 2 minutes.
4. Available currencies are refreshed only on cache miss and the expiry of that is 1 week

periodic jobs:
1. refresher
2. expired entries cleaner