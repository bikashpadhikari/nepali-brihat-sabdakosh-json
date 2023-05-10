# Structured JSON of Nepal Brihat Sabdakosh

This repository contains a structured JSON dump of all 122,000 words of the
Nepali Brihat Sabdakosh (नेपाली बृहत् शब्दकोश, *Unabridged Nepali Dictionary*)
published by Nepal Academy. It also contains the tools necessary to generate
the JSON.

## Data Source

Data was extracted from the version 19 APK of the `np.com.naya.sabdakosh`
Android app ([play store
link](https://play.google.com/store/apps/details?id=np.com.naya.sabdakosh)).
The app uses an embedded realm database encrypted with the following key:

```
7feb8c9cd654567106e867956802fd25609d58624539e10f02eaf6aef5facda9d5b9f5024fe4234c1f08c01ed875976719369dfa94b645a1212fdd968e00b6f3
```

You'll have to use [Realm Studio version
13](https://github.com/realm/realm-studio/releases/tag/v13.0.2) to open the
Realm database embedded in the app under `assets/db.realm`, since newer versions
aren't backwards compatible with the version the app uses. `extract.go`
assumes that the input was exported using Realm Studio's "export to
JSON" feature.

Each meaning in the Realm database is also encrypted with an AES-256-CBC
passphrase of `058aa5325d7d2e7`. You can decrypt individual rows with a command
like the following:

```
echo "U2FsdGVkX19pmyuNuE6X1Cne+Qc2mEhxBXrawcMdh/tkkZnj7Dj2Z0HYGPCQl27RM30pTEvYM6VuAK/WZtlJh07YkLaM6CRJI6XjrL4egaHF3ijpm/kuyT7hzQjHOU2gRtJNLCFXTbLP/RHUPj1+sHNylAmsbnI8zHSO7C
PU61A=" | openssl enc -aes-256-cbc -d -a -A -md md5 -pass pass:058aa5325d7d2e7
<span▥>चिकामारी</span><br/><br/><a◳>ना.</a><p▦>चिकीखेल।</p> 
```

## Data Schema

Each entry is a JSON object. For example, for अ, the object is:

```
{
  "word": "अ",
  "definitions": [
    {
      "grammar": "ना.",
      "senses": [
        "१. देवनागरी वर्णमालाको स्वर वर्णमध्ये पहिलो स्वर वर्ण; परम्परागत रूपमा कण्ठस्थानबाट उच्चारण हुने ह्रस्व स्वर वर्ण र भाषाविज्ञानअनुसार आधा खुला; केन्द्रीय स्वर वर्ण; लेख्य रूपमा सो स्वर वर्णको प्रतिनिधित्व गर्ने लिपिचिह्न।",
        "२. लेखाइका क्रममा विषयको विभाजन उपविभाजनका निम्ति स्वर वर्णको प्रयोग गरिँदा दिइने क्रमबोधक पहिलो चिह्न।"
      ]
    },
    {
      "grammar": "ना.",
      "etymology": "[सं.]",
      "senses": [
        "१. संस्कृत एकाक्षरी कोशअनुसार मूलतः विष्णुलाई जनाउने मङ्गलवाची शब्द।",
        "२. ॐ भित्र निहित अ+उ+म् तीन ध्वनिमध्ये विष्णुलाई बुझाउने पहिलो ध्वनि (उ तथा म् ध्वनि क्रमशः शिव तथा ब्रह्मालाई बुझाउने मानिन्छन्)।"
      ]
    },
    {
      "grammar": "नि.",
      "senses": [
        "झर्को, गाली, बेवास्ता, अस्वीकार आदि बुझाउन आवेगका अवस्थामा प्रयोग गरिने विस्मयादिबोधक शब्द; आ।"
      ]
    },
    {
      "grammar": "पूस.",
      "senses": [
        "शब्दका अगाडि लागेर अभाव, भिन्नता, विपरीतता आदि बुझाउने पूर्वसर्ग।"
      ]
    },
    {
      "grammar": "नि.",
      "senses": [
        "दिक्क लागेको अवस्थामा व्यक्त गरिने उपेक्षा भाव।"
      ]
    }
  ]
}
```

Each `definition` can have a `grammar`, `etymology` and `senses` field.
A `sense` is actually an HTML string, and may include examples tagged by
a `<span class="example">`.


## License

All code is under the MIT license. I'm not sure what the license for the actual
dictionary would be. Use at your own risk.
