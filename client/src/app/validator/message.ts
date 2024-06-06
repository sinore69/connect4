import Ajv, {JSONSchemaType} from "ajv"
const ajv = new Ajv()
interface Message{
    Text:string
}

const schema: JSONSchemaType<Message> = {
  type: "object",
  properties: {
    Text: {type: "string"},
  },
  required: ["Text"],
  additionalProperties: false
}

const validate = ajv.compile(schema)

export function isMessage(data:any){
    if (validate(data)) {
      return true
    } else {
      return false
    }
}