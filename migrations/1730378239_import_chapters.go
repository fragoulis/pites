package migrations

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"

	"github.com/fragoulis/setip_v2/internal/app/chapter/model"
)

//nolint:lll
func init() {
	m.Register(func(db dbx.Builder) error {
		chaptersRaw := `
[
 {
  "city_ids": [
   "jbtse3nne5pfr6e",
   "zr4owrcdcb1u3bp",
   "agydtplhz4qj75l",
   "m7nnvfnk674hoiz",
   "0jb44p7eob2bn8r",
   "4a77j9w8ft63zzt",
   "98qaun1rwnj97pv",
   "4tbymejvva6vjy2",
   "x32mio25ozgym5g",
   "57jlpb79p8to6pn",
   "j42unwxthqdeyzi",
   "x9c7uyajrniy369",
   "y2uu9in4dvi4t3o",
   "a5d9nk6of1t2qw4",
   "315zb5lbjwqvb9e",
   "8e556c1i6ihb2lu",
   "fv15odxn7cdgw7w",
   "w26sa7w90a9c6xs",
   "k3hlgfcpw51exrd",
   "peoilt1xq5e10kf",
   "ohbacpqns4vdup5",
   "kcp9wpeszweo5pj",
   "0w4jv81bk74qzys",
   "8jqfqh16mhx1i09"
  ],
  "created": "2024-10-31 07:35:37.755Z",
  "id": "r4snm4p9f78f6zq",
  "name": "ΒΟΡΕΙΑ",
  "raw_city_query": "ιωνια,παρασκευη,αγιος στεφανος,ανθουσα,ανοιξη,βιρλησσια,βριλησσια,γαλατσι,γερακα,δροσια,κηφισια,κοκκινος μυλος,κρυονερι,λυκοβρυση,μαδυτος,μαρουσ,μελισσια,μελλησια,μεταμορφωση,χαλκηδο,ερυθραια,ιωανια,ιωνια,ιωννια,ηρακλειο,παπαγο,πεντελη,περισσος,πευκακια,πευκη,ριζουπολη,φιλαδε,χαλανδρ,χολαργο,ψυχικο,παρασκευη",
  "updated": "2024-10-31 12:26:25.923Z"
 },
 {
  "city_ids": [
   "qagtx32zn20qxhk",
   "eadqz4anq1senho",
   "zmwocsav6a8cpsc",
   "k2cxewv2z2xtfye"
  ],
  "created": "2024-10-31 12:09:36.413Z",
  "id": "d4dgyu0ria4gt7p",
  "name": "ΑΝΑΤΟΛΙΚΑ",
  "raw_city_query": "ιλισια,βυρωνα,ζωγραφο,ιλισια,καισαριανη,καισιαρανη,καισιαριανη,υμηττο,υμμητο",
  "updated": "2024-10-31 12:09:36.413Z"
 },
 {
  "city_ids": [
   "bluj29lf6uxb92l",
   "qr1vel766ilx2dz",
   "tf9z9o9vrv4jxsy",
   "rb4dz0u127y82zs",
   "13kp04nseyvl25x",
   "lnnvn5l910oavm3",
   "ezwnfks6soa56ie",
   "pf7ldyj1z3x9fqd",
   "zrawnsiz6e43v9b"
  ],
  "created": "2024-10-31 12:12:26.900Z",
  "id": "z5smt1qucu8r2j1",
  "name": "ΠΕΙΡΑΙΑΣ",
  "raw_city_query": "αγ. ιωαννης,αγιος ιωαννης,αγια σοφια,ρεντη,δραπετσωνα,καλιπολη,καλλιπολη,καμινια,καστελλα,κερατσινι,κορυδαλ,κορυδαλ,νικαια,πειραια,περαμα,σαλαμινα,φαληρο",
  "updated": "2024-10-31 12:12:26.900Z"
 },
 {
  "city_ids": [
   "aenbfrf7vw57pom",
   "0q02xv4n5xa3dot",
   "zqocfgsznecwtqw",
   "2c9awwc7l0azxv1",
   "fkate42z75zbcg7",
   "sciaqj73la4pic7",
   "4iqf3r60l3pty6g",
   "xtoeo4fia7mfyn1",
   "8lgqcr3kfa6doyx",
   "v4dl1a1xvtgmfkh",
   "57z3i5l3dmhnljf",
   "zrawnsiz6e43v9b",
   "3y1gygri863ucig",
   "j9h7g8t48drgzbo",
   "ihuec5w7gclkcwe",
   "sbcc2asqgi6srwk",
   "f52v37xiegdrtz1",
   "c82w1v26n3rk5r7",
   "zpt390mj1xuyo7d",
   "iggjuda3bk8k2df",
   "1xb6sylkedxm83c",
   "4c7xx0psn2qdk34",
   "442h7xd5aai4vs0",
   "lh1qorlttr4csny",
   "uf2cfm5xfakyz99",
   "vvnq2o0ajs9crmt"
  ],
  "created": "2024-10-31 12:19:18.718Z",
  "id": "9evse27yrtdtf8l",
  "name": "ΝΟΤΙΑ",
  "raw_city_query": "α. γλυφαδα,αγ. δημητριος,αγ.δημητριος,αγιος  δημητριος,αγιος δημητριος,αλιμος,ανω βουλα,ανω γλυφαδα,αργυροπολη,αργυρουπολη,βουλα,γλυφαδα,δαφνη,ελληνικο,ηλιουπολη,καλλιθεα,μοσχατο,ν. σμυρνη,ν.σμυρνη,νεα σμυρνη,π. φαληρο,π.φαληρο,παλαιο φαληρο,ταυρος,τζιτζιφιες,δυτικα,βαρβαρα,αναργυρ,αγιοι αναργυ,αγιος φανουρ,αιγαλεω,ανθουπολη,λιοσι,απροπυργο,ασπροπυργο,αχαρναι,αχαρνες,ζεφυρι,ηλιον,ιλιον,καματερο,κηπουπολη,μενιδι,νεα ζωη,περιστερι,πετρουπολ,φυλη,χαιδαρι",
  "updated": "2024-10-31 12:19:18.718Z"
 },
 {
  "city_ids": [
   "r3c601mz97dq7ac"
  ],
  "created": "2024-10-31 12:31:17.859Z",
  "id": "u2pxk2hhubat759",
  "name": "ΚΕΝΤΡΟ",
  "raw_city_query": "κυψελη,αρτεμιος,ελευθεριος,παντελεημωνας,νικολαος,αθην,ακροπολη,αμπελοκηποι,πετραλω,πατησια,αττικη,βικτωρια,βοτανικος,γκαζι,γκυζη,γουδη,γουδι,ελληνορωσων,εξαρχεια,θησειο,θυμαρακια,κεντρο,κεραμεικος,κολωνο,κολωνο,κωλονο,κουκακι,κυπριαδου,λαμπρινη,σκουζε,στρεφη,λυκαββητ,λυκαβηττ,μετς,μοναστηρακι,ν. κοσμος,νεο κοσμο,νεος κοσμος,νεαπολη,παγκρατ,παγρατι,πανορμου,πατισια,βαθης,πλακα,αμερικης,πολυγωνο,σεπολια,χιλτον",
  "updated": "2024-10-31 12:31:17.859Z"
 }
]
`

		var chapters []*model.Chapter

		err := json.Unmarshal([]byte(chaptersRaw), &chapters)
		if err != nil {
			return err
		}

		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("u44zq22ph4yrurl")
		if err != nil {
			return err
		}

		return dao.RunInTransaction(func(tx *daos.Dao) error {
			for _, chapter := range chapters {
				record := models.NewRecord(collection)
				record.SetId(chapter.ID)
				record.Set("name", chapter.Name)
				record.Set("raw_city_query", chapter.RawCityQuery)
				record.Set("city_ids", chapter.CityIDs)

				if err := tx.SaveRecord(record); err != nil {
					return fmt.Errorf("failed to save record: %w", err)
				}
			}

			return nil
		})
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery("delete from chapters").Execute()
		if err != nil {
			return err
		}

		return nil
	})
}
