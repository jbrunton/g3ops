class Application
  def call(env)
    status  = 200
    headers = { "Content-Type" => "application/json" }
    body    = ['{ "greeting": "Hello, World!" }']

    [status, headers, body]
  end
end

run Application.new
