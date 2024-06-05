import Ajv, {JSONSchemaType} from "ajv"
const ajv = new Ajv()
interface InitialState{
    Disable:boolean
}

const schema: JSONSchemaType<InitialState> = {
  type: "object",
  properties: {
    Disable: {type: "boolean"},
  },
  required: ["Disable"],
  additionalProperties: false
}

const validate = ajv.compile(schema)

export function isInitialState(data:any){
    if (validate(data)) {
      return true
    } else {
      return false
    }
}