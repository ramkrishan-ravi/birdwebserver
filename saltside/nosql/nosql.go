package nosql
import(
)
import (
	"log"
	"fmt"
	"time"
	"github.com/bradfitz/gomemcache/memcache"
)

const(
	memcacheTimeOut = 300 //connection timeout to memcache in ms
	server = "127.0.0.0.1:11211" //server ip where memcache is running
	defTTL = 3600 //time in second for which data will stay in memcache after that will be deleted
)
type MemcacheClient struct{
	mc *memcache.Client
}

var MemClient MemcacheClient

//This function will initialize cache handle with required data like timeout ,server list to whom it need to connect etc
func init(){
	
	MemClient.mc = memcache.New(server)

	if MemClient.mc == nil {
		log.Panic("Error: unable to initialize memcache \n")
	}

	MemClient.mc.Timeout = time.Duration(memcacheTimeOut) * time.Duration(time.Millisecond)
}
//Method to get memcache handle one handle is capable of doing multiple concurrent operation
func GetHandle() *MemcacheClient{
	return &MemClient
}
/*Method to get data from memcache
input 
key
output
on success will return data fetched from memacache and error will be nil
on fail will return nil as data and error
*/
func ( handle * MemcacheClient) Get(key string) ([]byte , error) {
	
	it , err := handle.mc.Get(key)
	if err != nil || it == nil {
		log.Printf("Error: failed to fetch key:%s %s",key,err)
		return nil, err
	}
	log.Printf("Fetching value for %s key success value:%s\n",key,it.Value)
	return it.Value, nil
}
/*Method to delete key from memcache
input 
key
output
on success will return error as nil otherwise actual error
*/
func ( handle * MemcacheClient) Delete(key string ) error {

	if len(key) <= 0 {
		log.Printf("Error: key passed is blank blank\n")
		return fmt.Errorf("key Passed is blank")
	}
	
	err := handle.mc.Delete(key)
	if err != nil{
		log.Printf("Error: unable to delete key:%s:%s\n",key,err.Error())
		return err
	}
	log.Printf("Info:Key:%s deleted successfully\n",key)
	return nil
}
/*Method to add bird in  memcache
input 
key
output
on success will return error as nil otherwise actual error
*/
func ( handle *MemcacheClient) Add(key string, value []byte) error {

	if len(key) <= 0 {
		log.Printf("Error: passed key to be placd in memcache in blank\n")
		return fmt.Errorf("passed key to be placd in memcache in blank")
	}
	log.Printf("Info: memcache set key:%s value:%s\n",key , string(value))
	it := memcache.Item{Key:key, Value:value, Expiration: defTTL}
	err := handle.mc.Set(&it)
	if err != nil{
		log.Printf("Error: unable to set data:%s:%s\n",key,err.Error())
		return fmt.Errorf("unable to set data:%s\n",err.Error())
	}

	return nil
}

/*scanning full memcache is not supported by this go client otherwise it will be like fetch all and send array of all id*/
func ( handle * MemcacheClient) GetAllKey( ) ([]string ,error){
	return nil,nil
}

