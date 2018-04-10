import {Pipe, PipeTransform} from "@angular/core";

@Pipe({
  name: 'search',
  pure: false
})

export class SearchPipe implements PipeTransform{
  transform(value, args) {
    //args is search string, args[0] is the first char, value is object array contains each line in table
    //in old angualr2, args[0] is the whole search string, so need to use args instead of args[0] below
    if (!args) { //search string is empty
      return value;
    } 
    else if (value) //search string is not empty and table has rows.
    {
      return value.filter(
        item => {
        for (let key in item) 
        {
          if ((typeof item[key] === 'string' || item[key] instanceof String) &&
            (item[key].indexOf(args) !== -1)) {
            return true;
          }
        }
      });
    }
  }
}

@Pipe({ 
  name: 'orderBy', 
  pure: true
})

export class SortPipe implements PipeTransform {
  transform(array: any, orderField :any): Array<string> {
    if (!Array.isArray(array)) { //if it's not array, sort menthod can not be used
      return array;
    }
    var orderType: boolean = (orderField.substr(0, 1) == "+"); //ascending, arrow up, 1,2,3
    orderField = orderField.substr(1);
    array.sort((a: any, b: any) => {
      if (a[orderField] == null || a[orderField].isUndefined 
        || b[orderField] == null || b[orderField].isUndefined || a[orderField] == b[orderField])
          return 0;

      if (orderType){
        return a[orderField] < b[orderField] ? -1 : 1;
      }else{
        return a[orderField] > b[orderField] ? -1 : 1;
      }
    });

    return array;
  }
}
