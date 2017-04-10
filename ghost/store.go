// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ghost

import "strings"

type store struct{}

func (s store) executeQuery(q Query) (string, error) {

	var JSONResponse string

	//Caching case
	//Return the cached result if there is a cache key present
	//AND there is a result from the cache
	if q.cacheKey != "" {
		cacheResult, ok := App.QueryCache.Get(q.cacheKey)
		if ok {
			LogDebug("STORE", true, "Returning from cache, using key: "+q.cacheKey, nil)
			return cacheResult.(string), nil
		}
	}

	//No caching case
	if err := App.DB.QueryRow(q.SQL).Scan(&JSONResponse); err != nil {
		//Only one row is returned as JSON is returned by Postgres
		//Empty result
		if strings.Contains(err.Error(), "sql") {
			return "", nil
		}

		//Else its a database error
		return "", err

	}

	//Set the cache if a cache key has been provided
	if q.cacheKey != "" {
		LogDebug("STORE", true, "Caching result with key: "+q.cacheKey, nil)
		App.QueryCache.Set(q.cacheKey, JSONResponse)
	}
	return JSONResponse, nil

}
