FROM onsdigital/dp-go

WORKDIR /app/

COPY ./build/dp-dd-search-indexer .

ENTRYPOINT ./dp-dd-search-indexer
