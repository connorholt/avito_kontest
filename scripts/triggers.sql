create or replace function update_timestamp_column()
returns trigger as $$
    begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language 'plpgsql';


create or replace trigger update_updated_at_column before update
    on banners for each row execute function update_timestamp_column();