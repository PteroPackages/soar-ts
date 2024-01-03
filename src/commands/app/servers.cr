module Soar::Commands::App
  class Servers < Base
    def setup : Nil
      @name = "servers"

      add_command List.new
    end

    def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
      stdout.puts help_template
    end

    private class List < Base
      @filters = [] of String

      def setup : Nil
        @name = "list"

        add_option "name", type: :single
        add_option "desc", type: :single
        add_option "uuid", type: :single
        add_option "image", type: :single
        add_option "json"
      end

      def pre_run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        super

        @filters << "name" if options.has? "name"
        @filters << "description" if options.has? "desc"
        @filters << "uuid" if options.has? "uuid"
        @filters << "image" if options.has? "image"
      end

      def run(arguments : Cling::Arguments, options : Cling::Options) : Nil
        path = "/api/application/servers"
        unless @filters.empty?
          def_filter_params name: "name", desc: "description", uuid: "uuid", image: "image"
        end

        servers, meta = request get: path, as: Array(Models::Server)
        if servers.empty?
          unless options.has? "json"
            stdout << "Showing 0 results from page "
            stdout << meta.current_page << " of " << meta.total_pages << '\n'
            stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
            stdout << "\n┃ ".colorize.light_gray

            if @filters.empty?
              stdout << "no filters applied"
            else
              stdout << "filters: ".colorize.dark_gray << @filters.join(", ")
            end
            return
          end
        end

        if options.has? "json"
          servers.to_json stdout
          return
        end

        width = 2 + (Math.log(servers.last.id.to_f + 1) / Math.log(10)).ceil.to_i
        servers.each do |server|
          server.to_s(stdout, width)
          stdout << "\n\n"
        end

        stdout << "Showing " << meta.count << " results from page "
        stdout << meta.current_page << " of " << meta.total_pages << '\n'
        stdout << "\n┃ ".colorize.light_gray << "total: ".colorize.dark_gray << meta.total
        stdout << "\n┃ ".colorize.light_gray

        if @filters.empty?
          stdout << "no filters applied"
        else
          stdout << "filters: ".colorize.dark_gray << @filters.join(", ")
        end
      end
    end
  end
end