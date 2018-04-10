import { Message } from './model';
export function formatTime(time): string { //input: 2016-10-31T09:18:36+09:00, output: 2016-10-31 09:18:36
  var result: string = "";
  result = time.replace(/T/, " ");
  result = result.replace(/\+09:00/, "");
  result = result.replace(/Z/, "");
  return result;
}

interface Response {
    code: number;
    msg: string;
}

export function getErrMsg(err: any): string { //input: 2016-10-31T09:18:36+09:00, output: 2016-10-31 09:18:36
  switch (err.status) {
    case 0:
      return " Internal server error";
    case 400:
      // if (err.statusText != null) //to handle blob response
      //   return err.statusText;
      let res: Response = JSON.parse(err._body);
      if( err._body.includes("Duplicate") || err._body.includes("already exists")){
        return "Data already exists";
      }else if(err._body.includes("json:") || err._body.includes("timeout") || err._body.includes("invalid character")){
        return " Internal server error";
      }
      return res.msg;
    case 401:
      return " Authentication failed";
    case 403:
      return " Permission denied";
    default:
      return "";
  }
}

export function createMsg(type, action, model, object, reason? : string): Message{
  if(type == "ok")
    return new Message("positive", "OK!", action + " " + model + ": " + object, "check circle");

  if (type == "error"){
    if(reason != "")
      return new Message("negative", "Failed!", "Failed to " + action + " " + model + ": " + object + ". " + reason, "warning circle");
    
    return new Message("negative", "Failed!", "Failed to " + action + " " + model + ": " + object, "warning circle");
  }
}