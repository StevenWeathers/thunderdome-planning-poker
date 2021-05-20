/*
list copied from countryflags.io with the following console snippet

let countries = []
$('.item_country p').map((i, c) => {
    if (i === 0 || i % 2 === 0) {
        const country = {
            abbrev: $(c).text(),
            name: $(c).next('p').text()
        }
        countries.push(country)
    }
})
copy(countries)

*/
export const countryList = [
    {
        abbrev: 'AD',
        name: 'Andorra',
    },
    {
        abbrev: 'AE',
        name: 'United Arab Emirates',
    },
    {
        abbrev: 'AF',
        name: 'Afghanistan',
    },
    {
        abbrev: 'AG',
        name: 'Antigua and Barbuda',
    },
    {
        abbrev: 'AI',
        name: 'Anguilla',
    },
    {
        abbrev: 'AH',
        name: '',
    },
    {
        abbrev: 'AK',
        name: '',
    },
    {
        abbrev: 'AL',
        name: 'Albania',
    },
    {
        abbrev: 'AM',
        name: 'Armenia',
    },
    {
        abbrev: 'AN',
        name: 'Netherlands Antilles',
    },
    {
        abbrev: 'AO',
        name: 'Angola',
    },
    {
        abbrev: 'AQ',
        name: 'Antarctica',
    },
    {
        abbrev: 'AR',
        name: 'Argentina',
    },
    {
        abbrev: 'AS',
        name: 'American Samoa',
    },
    {
        abbrev: 'AT',
        name: 'Austria',
    },
    {
        abbrev: 'AU',
        name: 'Australia',
    },
    {
        abbrev: 'AW',
        name: 'Aruba',
    },
    {
        abbrev: 'AX',
        name: 'Åland Islands',
    },
    {
        abbrev: 'AZ',
        name: 'Azerbaijan',
    },
    {
        abbrev: 'BA',
        name: 'Bosnia and Herzegovina',
    },
    {
        abbrev: 'BB',
        name: 'Barbados',
    },
    {
        abbrev: 'BD',
        name: 'Bangladesh',
    },
    {
        abbrev: 'BE',
        name: 'Belgium',
    },
    {
        abbrev: 'BF',
        name: 'Burkina Faso',
    },
    {
        abbrev: 'BG',
        name: 'Bulgaria',
    },
    {
        abbrev: 'BH',
        name: 'Bahrain',
    },
    {
        abbrev: 'BI',
        name: 'Burundi',
    },
    {
        abbrev: 'BJ',
        name: 'Benin',
    },
    {
        abbrev: 'BL',
        name: 'Saint Barthélemy',
    },
    {
        abbrev: 'BM',
        name: 'Bermuda',
    },
    {
        abbrev: 'BN',
        name: 'Brunei Darussalam',
    },
    {
        abbrev: 'BO',
        name: 'Bolivia',
    },
    {
        abbrev: 'BQ',
        name: 'Bonaire, Sint Eustatius and Saba',
    },
    {
        abbrev: 'BR',
        name: 'Brazil',
    },
    {
        abbrev: 'BS',
        name: 'Bahamas',
    },
    {
        abbrev: 'BT',
        name: 'Bhutan',
    },
    {
        abbrev: 'BV',
        name: 'Bouvet Island',
    },
    {
        abbrev: 'BW',
        name: 'Botswana',
    },
    {
        abbrev: 'BY',
        name: 'Belarus',
    },
    {
        abbrev: 'BZ',
        name: 'Belize',
    },
    {
        abbrev: 'CA',
        name: 'Canada',
    },
    {
        abbrev: 'CC',
        name: 'Cocos (Keeling) Islands',
    },
    {
        abbrev: 'CD',
        name: 'Congo, The Democratic Republic Of The',
    },
    {
        abbrev: 'CF',
        name: 'Central African Republic',
    },
    {
        abbrev: 'CG',
        name: 'Congo',
    },
    {
        abbrev: 'CH',
        name: 'Switzerland',
    },
    {
        abbrev: 'CI',
        name: "Côte D'Ivoire",
    },
    {
        abbrev: 'CK',
        name: 'Cook Islands',
    },
    {
        abbrev: 'CL',
        name: 'Chile',
    },
    {
        abbrev: 'CM',
        name: 'Cameroon',
    },
    {
        abbrev: 'CN',
        name: 'China',
    },
    {
        abbrev: 'CO',
        name: 'Colombia',
    },
    {
        abbrev: 'CR',
        name: 'Costa Rica',
    },
    {
        abbrev: 'CU',
        name: 'Cuba',
    },
    {
        abbrev: 'CV',
        name: 'Cape Verde',
    },
    {
        abbrev: 'CW',
        name: 'Curaçao',
    },
    {
        abbrev: 'CX',
        name: 'Christmas Island',
    },
    {
        abbrev: 'CY',
        name: 'Cyprus',
    },
    {
        abbrev: 'CZ',
        name: 'Czech Republic',
    },
    {
        abbrev: 'DE',
        name: 'Germany',
    },
    {
        abbrev: 'DJ',
        name: 'Djibouti',
    },
    {
        abbrev: 'DK',
        name: 'Denmark',
    },
    {
        abbrev: 'DM',
        name: 'Dominica',
    },
    {
        abbrev: 'DO',
        name: 'Dominican Republic',
    },
    {
        abbrev: 'DZ',
        name: 'Algeria',
    },
    {
        abbrev: 'EC',
        name: 'Ecuador',
    },
    {
        abbrev: 'EE',
        name: 'Estonia',
    },
    {
        abbrev: 'EG',
        name: 'Egypt',
    },
    {
        abbrev: 'EH',
        name: 'Western Sahara',
    },
    {
        abbrev: 'ER',
        name: 'Eritrea',
    },
    {
        abbrev: 'ES',
        name: 'Spain',
    },
    {
        abbrev: 'ET',
        name: 'Ethiopia',
    },
    {
        abbrev: 'EU',
        name: '',
    },
    {
        abbrev: 'FI',
        name: 'Finland',
    },
    {
        abbrev: 'FJ',
        name: 'Fiji',
    },
    {
        abbrev: 'FK',
        name: 'Falkland Islands (Malvinas)',
    },
    {
        abbrev: 'FM',
        name: 'Micronesia, Federated States Of',
    },
    {
        abbrev: 'FO',
        name: 'Faroe Islands',
    },
    {
        abbrev: 'FR',
        name: 'France',
    },
    {
        abbrev: 'GA',
        name: 'Gabon',
    },
    {
        abbrev: 'GB',
        name: 'United Kingdom',
    },
    {
        abbrev: 'GD',
        name: 'Grenada',
    },
    {
        abbrev: 'GE',
        name: 'Georgia',
    },
    {
        abbrev: 'GF',
        name: 'French Guiana',
    },
    {
        abbrev: 'GG',
        name: 'Guernsey',
    },
    {
        abbrev: 'GH',
        name: 'Ghana',
    },
    {
        abbrev: 'GI',
        name: 'Gibraltar',
    },
    {
        abbrev: 'GL',
        name: 'Greenland',
    },
    {
        abbrev: 'GM',
        name: 'Gambia',
    },
    {
        abbrev: 'GN',
        name: 'Guinea',
    },
    {
        abbrev: 'GP',
        name: 'Guadeloupe',
    },
    {
        abbrev: 'GQ',
        name: 'Equatorial Guinea',
    },
    {
        abbrev: 'GR',
        name: 'Greece',
    },
    {
        abbrev: 'GS',
        name: 'South Georgia and the South Sandwich Islands',
    },
    {
        abbrev: 'GT',
        name: 'Guatemala',
    },
    {
        abbrev: 'GU',
        name: 'Guam',
    },
    {
        abbrev: 'GW',
        name: 'Guinea-Bissau',
    },
    {
        abbrev: 'GY',
        name: 'Guyana',
    },
    {
        abbrev: 'HK',
        name: 'Hong Kong',
    },
    {
        abbrev: 'HM',
        name: 'Heard and McDonald Islands',
    },
    {
        abbrev: 'HN',
        name: 'Honduras',
    },
    {
        abbrev: 'HR',
        name: 'Croatia',
    },
    {
        abbrev: 'HT',
        name: 'Haiti',
    },
    {
        abbrev: 'HU',
        name: 'Hungary',
    },
    {
        abbrev: 'IC',
        name: '',
    },
    {
        abbrev: 'ID',
        name: 'Indonesia',
    },
    {
        abbrev: 'IE',
        name: 'Ireland',
    },
    {
        abbrev: 'IL',
        name: 'Israel',
    },
    {
        abbrev: 'IM',
        name: 'Isle of Man',
    },
    {
        abbrev: 'IN',
        name: 'India',
    },
    {
        abbrev: 'IO',
        name: 'British Indian Ocean Territory',
    },
    {
        abbrev: 'IQ',
        name: 'Iraq',
    },
    {
        abbrev: 'IR',
        name: 'Iran, Islamic Republic Of',
    },
    {
        abbrev: 'IS',
        name: 'Iceland',
    },
    {
        abbrev: 'IT',
        name: 'Italy',
    },
    {
        abbrev: 'JE',
        name: 'Jersey',
    },
    {
        abbrev: 'JM',
        name: 'Jamaica',
    },
    {
        abbrev: 'JO',
        name: 'Jordan',
    },
    {
        abbrev: 'JP',
        name: 'Japan',
    },
    {
        abbrev: 'KE',
        name: 'Kenya',
    },
    {
        abbrev: 'KG',
        name: 'Kyrgyzstan',
    },
    {
        abbrev: 'KH',
        name: 'Cambodia',
    },
    {
        abbrev: 'KI',
        name: 'Kiribati',
    },
    {
        abbrev: 'KM',
        name: 'Comoros',
    },
    {
        abbrev: 'KN',
        name: 'Saint Kitts And Nevis',
    },
    {
        abbrev: 'KP',
        name: "Korea, Democratic People's Republic Of",
    },
    {
        abbrev: 'KR',
        name: 'Korea, Republic of',
    },
    {
        abbrev: 'KW',
        name: 'Kuwait',
    },
    {
        abbrev: 'KY',
        name: 'Cayman Islands',
    },
    {
        abbrev: 'KZ',
        name: 'Kazakhstan',
    },
    {
        abbrev: 'LA',
        name: "Lao People's Democratic Republic",
    },
    {
        abbrev: 'LB',
        name: 'Lebanon',
    },
    {
        abbrev: 'LC',
        name: 'Saint Lucia',
    },
    {
        abbrev: 'LI',
        name: 'Liechtenstein',
    },
    {
        abbrev: 'LK',
        name: 'Sri Lanka',
    },
    {
        abbrev: 'LR',
        name: 'Liberia',
    },
    {
        abbrev: 'LS',
        name: 'Lesotho',
    },
    {
        abbrev: 'LT',
        name: 'Lithuania',
    },
    {
        abbrev: 'LU',
        name: 'Luxembourg',
    },
    {
        abbrev: 'LV',
        name: 'Latvia',
    },
    {
        abbrev: 'LY',
        name: 'Libya',
    },
    {
        abbrev: 'MA',
        name: 'Morocco',
    },
    {
        abbrev: 'MC',
        name: 'Monaco',
    },
    {
        abbrev: 'MD',
        name: 'Moldova, Republic of',
    },
    {
        abbrev: 'ME',
        name: 'Montenegro',
    },
    {
        abbrev: 'MF',
        name: 'Saint Martin',
    },
    {
        abbrev: 'MG',
        name: 'Madagascar',
    },
    {
        abbrev: 'MH',
        name: 'Marshall Islands',
    },
    {
        abbrev: 'MK',
        name: 'Macedonia, the Former Yugoslav Republic Of',
    },
    {
        abbrev: 'ML',
        name: 'Mali',
    },
    {
        abbrev: 'MM',
        name: 'Myanmar',
    },
    {
        abbrev: 'MN',
        name: 'Mongolia',
    },
    {
        abbrev: 'MO',
        name: 'Macao',
    },
    {
        abbrev: 'MP',
        name: 'Northern Mariana Islands',
    },
    {
        abbrev: 'MQ',
        name: 'Martinique',
    },
    {
        abbrev: 'MR',
        name: 'Mauritania',
    },
    {
        abbrev: 'MS',
        name: 'Montserrat',
    },
    {
        abbrev: 'MT',
        name: 'Malta',
    },
    {
        abbrev: 'MU',
        name: 'Mauritius',
    },
    {
        abbrev: 'MV',
        name: 'Maldives',
    },
    {
        abbrev: 'MW',
        name: 'Malawi',
    },
    {
        abbrev: 'MX',
        name: 'Mexico',
    },
    {
        abbrev: 'MY',
        name: 'Malaysia',
    },
    {
        abbrev: 'MZ',
        name: 'Mozambique',
    },
    {
        abbrev: 'NA',
        name: 'Namibia',
    },
    {
        abbrev: 'NC',
        name: 'New Caledonia',
    },
    {
        abbrev: 'NE',
        name: 'Niger',
    },
    {
        abbrev: 'NF',
        name: 'Norfolk Island',
    },
    {
        abbrev: 'NG',
        name: 'Nigeria',
    },
    {
        abbrev: 'NI',
        name: 'Nicaragua',
    },
    {
        abbrev: 'NL',
        name: 'Netherlands',
    },
    {
        abbrev: 'NO',
        name: 'Norway',
    },
    {
        abbrev: 'NP',
        name: 'Nepal',
    },
    {
        abbrev: 'NR',
        name: 'Nauru',
    },
    {
        abbrev: 'NU',
        name: 'Niue',
    },
    {
        abbrev: 'NY',
        name: '',
    },
    {
        abbrev: 'NZ',
        name: 'New Zealand',
    },
    {
        abbrev: 'OM',
        name: 'Oman',
    },
    {
        abbrev: 'PA',
        name: 'Panama',
    },
    {
        abbrev: 'PE',
        name: 'Peru',
    },
    {
        abbrev: 'PF',
        name: 'French Polynesia',
    },
    {
        abbrev: 'PG',
        name: 'Papua New Guinea',
    },
    {
        abbrev: 'PH',
        name: 'Philippines',
    },
    {
        abbrev: 'PK',
        name: 'Pakistan',
    },
    {
        abbrev: 'PL',
        name: 'Poland',
    },
    {
        abbrev: 'PM',
        name: 'Saint Pierre And Miquelon',
    },
    {
        abbrev: 'PN',
        name: 'Pitcairn',
    },
    {
        abbrev: 'PR',
        name: 'Puerto Rico',
    },
    {
        abbrev: 'PS',
        name: 'Palestine, State of',
    },
    {
        abbrev: 'PT',
        name: 'Portugal',
    },
    {
        abbrev: 'PW',
        name: 'Palau',
    },
    {
        abbrev: 'PY',
        name: 'Paraguay',
    },
    {
        abbrev: 'QA',
        name: 'Qatar',
    },
    {
        abbrev: 'RE',
        name: 'Réunion',
    },
    {
        abbrev: 'RO',
        name: 'Romania',
    },
    {
        abbrev: 'RS',
        name: 'Serbia',
    },
    {
        abbrev: 'RU',
        name: 'Russian Federation',
    },
    {
        abbrev: 'RW',
        name: 'Rwanda',
    },
    {
        abbrev: 'SA',
        name: 'Saudi Arabia',
    },
    {
        abbrev: 'SB',
        name: 'Solomon Islands',
    },
    {
        abbrev: 'SC',
        name: 'Seychelles',
    },
    {
        abbrev: 'SD',
        name: 'Sudan',
    },
    {
        abbrev: 'SE',
        name: 'Sweden',
    },
    {
        abbrev: 'SG',
        name: 'Singapore',
    },
    {
        abbrev: 'SH',
        name: 'Saint Helena',
    },
    {
        abbrev: 'SI',
        name: 'Slovenia',
    },
    {
        abbrev: 'SJ',
        name: 'Svalbard And Jan Mayen',
    },
    {
        abbrev: 'SK',
        name: 'Slovakia',
    },
    {
        abbrev: 'SL',
        name: 'Sierra Leone',
    },
    {
        abbrev: 'SM',
        name: 'San Marino',
    },
    {
        abbrev: 'SN',
        name: 'Senegal',
    },
    {
        abbrev: 'SO',
        name: 'Somalia',
    },
    {
        abbrev: 'SR',
        name: 'Suriname',
    },
    {
        abbrev: 'SS',
        name: 'South Sudan',
    },
    {
        abbrev: 'ST',
        name: 'Sao Tome and Principe',
    },
    {
        abbrev: 'SV',
        name: 'El Salvador',
    },
    {
        abbrev: 'SX',
        name: 'Sint Maarten',
    },
    {
        abbrev: 'SY',
        name: 'Syrian Arab Republic',
    },
    {
        abbrev: 'SZ',
        name: 'Swaziland',
    },
    {
        abbrev: 'TC',
        name: 'Turks and Caicos Islands',
    },
    {
        abbrev: 'TD',
        name: 'Chad',
    },
    {
        abbrev: 'TF',
        name: 'French Southern Territories',
    },
    {
        abbrev: 'TG',
        name: 'Togo',
    },
    {
        abbrev: 'TH',
        name: 'Thailand',
    },
    {
        abbrev: 'TJ',
        name: 'Tajikistan',
    },
    {
        abbrev: 'TK',
        name: 'Tokelau',
    },
    {
        abbrev: 'TL',
        name: 'Timor-Leste',
    },
    {
        abbrev: 'TM',
        name: 'Turkmenistan',
    },
    {
        abbrev: 'TN',
        name: 'Tunisia',
    },
    {
        abbrev: 'TO',
        name: 'Tonga',
    },
    {
        abbrev: 'TR',
        name: 'Turkey',
    },
    {
        abbrev: 'TT',
        name: 'Trinidad and Tobago',
    },
    {
        abbrev: 'TV',
        name: 'Tuvalu',
    },
    {
        abbrev: 'TW',
        name: 'Taiwan, Republic Of China',
    },
    {
        abbrev: 'TZ',
        name: 'Tanzania, United Republic of',
    },
    {
        abbrev: 'UA',
        name: 'Ukraine',
    },
    {
        abbrev: 'UG',
        name: 'Uganda',
    },
    {
        abbrev: 'UM',
        name: 'United States Minor Outlying Islands',
    },
    {
        abbrev: 'US',
        name: 'United States',
    },
    {
        abbrev: 'UY',
        name: 'Uruguay',
    },
    {
        abbrev: 'UZ',
        name: 'Uzbekistan',
    },
    {
        abbrev: 'VA',
        name: 'Holy See (Vatican City State)',
    },
    {
        abbrev: 'VC',
        name: 'Saint Vincent And The Grenadines',
    },
    {
        abbrev: 'VE',
        name: 'Venezuela, Bolivarian Republic of',
    },
    {
        abbrev: 'VG',
        name: 'Virgin Islands, British',
    },
    {
        abbrev: 'VI',
        name: 'Virgin Islands, U.S.',
    },
    {
        abbrev: 'VN',
        name: 'Vietnam',
    },
    {
        abbrev: 'VU',
        name: 'Vanuatu',
    },
    {
        abbrev: 'WF',
        name: 'Wallis and Futuna',
    },
    {
        abbrev: 'WS',
        name: 'Samoa',
    },
    {
        abbrev: 'XK',
        name: '',
    },
    {
        abbrev: 'YE',
        name: 'Yemen',
    },
    {
        abbrev: 'YT',
        name: 'Mayotte',
    },
    {
        abbrev: 'ZA',
        name: 'South Africa',
    },
    {
        abbrev: 'ZM',
        name: 'Zambia',
    },
    {
        abbrev: 'ZW',
        name: 'Zimbabwe',
    },
]

export const countryMap = countryList.reduce((prev, country) => {
    prev[country.abbrev] = country.name

    return prev
}, {})
