FROM kong:3.6

# Copy in your DB-less config
COPY kong.yml /usr/local/kong/declarative/kong.yml

# Tell Kong to run DB-less with that file
ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/usr/local/kong/declarative/kong.yml

# (Optional) expose the port you’ve configured in kong.yml
EXPOSE 8000
