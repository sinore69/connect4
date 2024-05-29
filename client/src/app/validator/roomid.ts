import Ajv, {JSONSchemaType} from "ajv"
const ajv = new Ajv()

interface MyData {
  Id: number
}

const schema: JSONSchemaType<MyData> = {
  type: "object",
  properties: {
    Id: {type: "integer"},
  },
  required: ["Id"],
  additionalProperties: false
}

// validate is a type guard for MyData - type is inferred from schema type
const validate = ajv.compile(schema)

// or, if you did not use type annotation for the schema,
// type parameter can be used to make it type guard:
// const validate = ajv.compile<MyData>(schema)

export function roomIdValidator(data:any){
    if (validate(data)) {
      // data is MyData here
      return true
    } else {
      console.log(validate.errors)
      return false
    }
}