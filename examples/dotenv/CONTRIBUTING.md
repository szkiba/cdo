# Example: dotenv

This example shows how to use the `.env` file. The following lines are included in the [`.env`](.env) file in this directory:

```env file=.env
answer=42
question="Doesn't matter"
```

## show - Print variable values 

This task displays some environment variables. 

```sh
echo "The answer: $answer"
echo "The question: ${question:-Any}"
```

Variable values ​​in the .env file can be overwritten using the `-e/--env` flag (`cdo show -e answer="Hello World!"`) or with the name=value parameter (`cdo show answer="Hello World!"`)

The variable values ​​specified in the [`.env`](.env) file can also be overwritten in the `.env.local` file. Since `.env.local` is conveniently included in [`.gitignore`](.gitignore), it can be used for local settings.
