package bird

import (
	"encoding/json"
	"time"
	"log"
)
/*
structure to hold bird data
used in string formation from structure and string from structure
*/
type Bird struct{
	ID       * string    `json:"id"`
	Name     * string    `json:"name"`
	Family   * string    `json:"family"`
	Added    * string    `json:"added"`
	Visible  * bool      `json:"visible"`
	Continents []string  `json:"continents"`
}
/*
structre used to hold array of all bird id .
will be used for string to json conversion and vice versa
*/
type BirdList struct{
	id []string
}
/*
Methid will return true if mandatory param is present else false
*/
func (bird * Bird)IsValid()bool{
	if bird.ID == nil || bird.Name == nil ||bird.Family == nil ||  bird.Continents == nil{
		return false
	}
	return true
}
/*this method will set def value for bird*/
func (bird * Bird) SetDefault(){
	if bird.Visible == nil {
		bird.Visible = new(bool)
		*bird.Visible = false
	}

	if bird.Added == nil {
		bird.Added = new(string)
		*bird.Added = time.Now().UTC().Format(time.UnixDate)
	}
}
/*
Method to parse string of byte and putting data in Bird structure
*/
func Parse(data []byte,bird *Bird) error {
	log.Printf("json is %s\n",data)
	if err := json.Unmarshal(data, bird); err != nil {
		log.Printf("%s\n",err)
		return err
	}
	log.Printf("id:%s ",*bird.ID)
	if bird.Name != nil {log.Printf("name:%s\n",*bird.Name)}
	if bird.Family != nil {log.Printf("family:%s\n",*bird.Family)}
	if bird.Visible != nil {log.Printf("visible:%d\n",*bird.Visible)}
	if bird.Continents != nil{log.Printf("cont %s\n",bird.Continents)}
	
	return  nil
}
