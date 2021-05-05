package source_test

import (
	"bytes"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUspBr_Transform(t *testing.T) {
	s := source.NewUspBr()

	events, err := s.Transform(bytes.NewReader([]byte(uspBrResponse)))
	if err != nil {
		t.Errorf("cannot transform events: %v", err)
		return
	}

	assert.Equal(t, 2, len(events))
	eventsValid(t, events)
}

var uspBrResponse = `
<q:quakeml xmlns="http://quakeml.org/xmlns/bed/1.2" xmlns:q="http://quakeml.org/xmlns/quakeml/1.2">
<script/>
<eventParameters publicID="smi:org.gfz-potsdam.de/geofon/EventParameters">
<event publicID="smi:org.gfz-potsdam.de/geofon/usp2021ilss">
<description>
<text>Panama</text>
<type>region name</type>
</description>
<typeCertainty>known</typeCertainty>
<creationInfo>
<agencyID>USP</agencyID>
<author>scevent@seisMaster</author>
<creationTime>2021-05-01T06:20:05.29143Z</creationTime>
</creationInfo>
<magnitude publicID="smi:org.gfz-potsdam.de/geofon/Magnitude/20210501133608.267671.26275">
<stationCount>28</stationCount>
<creationInfo>
<agencyID>USP</agencyID>
<author>jroberto</author>
<creationTime>2021-05-01T13:36:08.267755Z</creationTime>
</creationInfo>
<mag>
<value>5.134310147</value>
<uncertainty>0.1846860909</uncertainty>
</mag>
<type>mb</type>
<originID>smi:org.gfz-potsdam.de/geofon/Origin/20210501133546.8097.26256</originID>
<methodID>smi:org.gfz-potsdam.de/geofon/trimmed_mean</methodID>
</magnitude>
<origin publicID="smi:org.gfz-potsdam.de/geofon/Origin/20210501133546.8097.26256">
<time>
<value>2021-05-01T06:13:39.256817Z</value>
<uncertainty>1.175634265</uncertainty>
</time>
<longitude>
<value>-79.20690155</value>
<uncertainty>5.983501434</uncertainty>
</longitude>
<latitude>
<value>9.512276649</value>
<uncertainty>5.075283527</uncertainty>
</latitude>
<quality>
<associatedPhaseCount>38</associatedPhaseCount>
<usedPhaseCount>38</usedPhaseCount>
<associatedStationCount>38</associatedStationCount>
<usedStationCount>38</usedStationCount>
<depthPhaseCount>0</depthPhaseCount>
<standardError>0.6823476109</standardError>
<azimuthalGap>104.6664734</azimuthalGap>
<maximumDistance>151.7942963</maximumDistance>
<minimumDistance>4.858849525</minimumDistance>
<medianDistance>36.93169022</medianDistance>
</quality>
<evaluationMode>manual</evaluationMode>
<creationInfo>
<agencyID>USP</agencyID>
<author>jroberto</author>
<creationTime>2021-05-01T13:35:46.810308Z</creationTime>
</creationInfo>
<depth>
<value>39298.46954</value>
<uncertainty>11788.76209</uncertainty>
</depth>
<originUncertainty>
<minHorizontalUncertainty>7698.101044</minHorizontalUncertainty>
<maxHorizontalUncertainty>14948.75908</maxHorizontalUncertainty>
<azimuthMaxHorizontalUncertainty>53.15739059</azimuthMaxHorizontalUncertainty>
<preferredDescription>horizontal uncertainty</preferredDescription>
</originUncertainty>
<methodID>smi:org.gfz-potsdam.de/geofon/LOCSAT</methodID>
<earthModelID>smi:org.gfz-potsdam.de/geofon/iasp91</earthModelID>
<evaluationStatus>confirmed</evaluationStatus>
</origin>
<preferredOriginID>smi:org.gfz-potsdam.de/geofon/Origin/20210501133546.8097.26256</preferredOriginID>
<preferredMagnitudeID>smi:org.gfz-potsdam.de/geofon/Magnitude/20210501133608.267671.26275</preferredMagnitudeID>
<type>earthquake</type>
</event>
<event publicID="smi:org.gfz-potsdam.de/geofon/usp2021iljg">
<description>
<text>Near East Coast of Honshu, Japan</text>
<type>region name</type>
</description>
<typeCertainty>known</typeCertainty>
<creationInfo>
...
</creationInfo>
<magnitude publicID="smi:org.gfz-potsdam.de/geofon/Magnitude/20210501132005.889567.20941">
<stationCount>6</stationCount>
<creationInfo>
<agencyID>USP</agencyID>
<author>jroberto</author>
<creationTime>2021-05-01T13:20:05.889586Z</creationTime>
</creationInfo>
<mag>
<value>6.845128063</value>
<uncertainty>0.1625816394</uncertainty>
</mag>
<type>mB</type>
<originID>smi:org.gfz-potsdam.de/geofon/Origin/20210501131922.91643.20883</originID>
<methodID>smi:org.gfz-potsdam.de/geofon/trimmed_mean</methodID>
<evaluationStatus>confirmed</evaluationStatus>
</magnitude>
<origin publicID="smi:org.gfz-potsdam.de/geofon/Origin/20210501131922.91643.20883">
<time>
<value>2021-05-01T01:27:21.072171Z</value>
<uncertainty>0.1512223482</uncertainty>
</time>
<longitude>
<value>141.8555298</value>
<uncertainty>5.338309288</uncertainty>
</longitude>
<latitude>
<value>38.11962891</value>
<uncertainty>5.278120041</uncertainty>
</latitude>
<depthType>operator assigned</depthType>
<quality>
<associatedPhaseCount>57</associatedPhaseCount>
<usedPhaseCount>55</usedPhaseCount>
<associatedStationCount>57</associatedStationCount>
<usedStationCount>55</usedStationCount>
<depthPhaseCount>0</depthPhaseCount>
<standardError>0.9511866921</standardError>
<azimuthalGap>52.32162476</azimuthalGap>
<maximumDistance>164.1687775</maximumDistance>
<minimumDistance>3.30714345</minimumDistance>
<medianDistance>148.1644592</medianDistance>
</quality>
<evaluationMode>manual</evaluationMode>
<creationInfo>
<agencyID>USP</agencyID>
<author>jroberto</author>
<creationTime>2021-05-01T13:19:22.917237Z</creationTime>
</creationInfo>
<depth>
<value>10000</value>
<uncertainty>0</uncertainty>
</depth>
<originUncertainty>
<minHorizontalUncertainty>9304.222107</minHorizontalUncertainty>
<maxHorizontalUncertainty>13138.41629</maxHorizontalUncertainty>
<azimuthMaxHorizontalUncertainty>134.0214081</azimuthMaxHorizontalUncertainty>
<preferredDescription>horizontal uncertainty</preferredDescription>
</originUncertainty>
<methodID>smi:org.gfz-potsdam.de/geofon/LOCSAT</methodID>
<earthModelID>smi:org.gfz-potsdam.de/geofon/iasp91</earthModelID>
<evaluationStatus>confirmed</evaluationStatus>
</origin>
<preferredOriginID>smi:org.gfz-potsdam.de/geofon/Origin/20210501131922.91643.20883</preferredOriginID>
<preferredMagnitudeID>smi:org.gfz-potsdam.de/geofon/Magnitude/20210501132005.889567.20941</preferredMagnitudeID>
<type>earthquake</type>
</event>
</eventParameters>
</q:quakeml>
`
