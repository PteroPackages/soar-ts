module Soar::Commands::App
  abstract class Base < Soar::Commands::Base
    protected getter config : Soar::Config { raise "unreachable" }
    protected getter client : Crest::Resource { raise "unreachable" }

    def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Bool
      return false unless super

      @config = Soar::Config.load_with options
      @client = Crest::Resource.new(
        config.app.url,
        headers: {
          "Authorization" => "Bearer #{config.app.key}",
          "Content-Type"  => "application/vnd.pterodactyl.v1+json",
          "Accept"        => "application/vnd.pterodactyl.v1+json",
        })

      true
    end

    protected def request(*, get path : String, as type : Array(T).class) : {Array(T), Models::FractalMeta} forall T
      res = client.get path
      data = Models::FractalList(T).from_json(res.body)

      {data.data.map(&.attributes), data.meta}
    end

    protected def request(*, get path : String, as type : T.class) : T forall T
      res = client.get path
      Models::FractalItem(T).from_json(res.body).attributes
    end

    protected def request(*, post path : String, data : _, as type : T.class) : T forall T
      res = client.post path, data
      Models::FractalItem(T).from_json(res.body).attributes
    end

    protected def request(*, delete path : String) : Nil
      client.delete path
    end

    private macro def_filter_params(**options)
      path += "?"
      path += URI::Params.build do |params|
        {% for key, name in options %}
          if options.has? {{ key.stringify }}
            params.add "filter[{{ name.id }}]", options.get({{ key.stringify }}).as_s
          end
        {% end %}
      end
    end
  end
end
