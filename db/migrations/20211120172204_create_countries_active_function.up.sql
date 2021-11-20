-- Get a list of countries
CREATE OR REPLACE FUNCTION countries_active() RETURNS table (
    country VARCHAR(2)
) AS $$
BEGIN
    RETURN QUERY SELECT ac.country FROM active_countries ac;
END;
$$ LANGUAGE plpgsql;