package source_test

import (
	"bytes"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeofon_Transform(t *testing.T) {
	s := source.NewGeofon()

	events, err := s.Transform(bytes.NewReader([]byte(geofonResponse)))
	if err != nil {
		t.Errorf("cannot transform events: %v", err)
		return
	}

	assert.Equal(t, 2, len(events))
	eventsValid(t, events)
}

var geofonResponse = `
<q:quakeml xmlns="http://quakeml.org/xmlns/bed/1.2" xmlns:q="http://quakeml.org/xmlns/quakeml/1.2">
<script/>
<eventParameters publicID="smi:org.gfz-potsdam.de/geofon/EventParameters">
<event publicID="smi:org.gfz-potsdam.de/geofon/gfz2021inzw">
<description>
<text>Southwest of Sumatra, Indonesia</text>
<type>region name</type>
</description>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T12:13:14.600338Z</creationTime>
</creationInfo>
<magnitude publicID="smi:org.gfz-potsdam.de/geofon/Origin/20210502134934.406947.1961745/netMag/mb">
<stationCount>43</stationCount>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T13:49:34.526352Z</creationTime>
</creationInfo>
<mag>
<value>4.767178353</value>
<uncertainty>0.2772048796</uncertainty>
</mag>
<type>mb</type>
<originID>smi:org.gfz-potsdam.de/geofon/Origin/20210502134934.406947.1961745</originID>
<methodID>smi:org.gfz-potsdam.de/geofon/trimmed_mean(25)</methodID>
</magnitude>
<origin publicID="smi:org.gfz-potsdam.de/geofon/Origin/20210502134934.406947.1961745">
<time>
<value>2021-05-02T12:07:30.939446Z</value>
<uncertainty>0.1781221777</uncertainty>
</time>
<longitude>
<value>103.3460922</value>
<uncertainty>2.12136507</uncertainty>
</longitude>
<latitude>
<value>-6.147405148</value>
<uncertainty>2.612462759</uncertainty>
</latitude>
<quality>
<associatedPhaseCount>78</associatedPhaseCount>
<usedPhaseCount>70</usedPhaseCount>
<associatedStationCount>76</associatedStationCount>
<usedStationCount>68</usedStationCount>
<standardError>0.9909751132</standardError>
<azimuthalGap>79.26596069</azimuthalGap>
<maximumDistance>88.75205994</maximumDistance>
<minimumDistance>1.301698089</minimumDistance>
<medianDistance>21.2748127</medianDistance>
</quality>
<evaluationMode>automatic</evaluationMode>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T13:49:34.407997Z</creationTime>
</creationInfo>
<depth>
<value>36038.20038</value>
<uncertainty>0</uncertainty>
</depth>
<methodID>smi:org.gfz-potsdam.de/geofon/LOCSAT</methodID>
<earthModelID>smi:org.gfz-potsdam.de/geofon/iasp91</earthModelID>
</origin>
<preferredOriginID>smi:org.gfz-potsdam.de/geofon/Origin/20210502134934.406947.1961745</preferredOriginID>
<preferredMagnitudeID>smi:org.gfz-potsdam.de/geofon/Origin/20210502134934.406947.1961745/netMag/mb</preferredMagnitudeID>
</event>
<event publicID="smi:org.gfz-potsdam.de/geofon/gfz2021inyh">
<description>
<text>Off Coast of Pakistan</text>
<type>region name</type>
</description>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T13:18:01.70721Z</creationTime>
</creationInfo>
<magnitude publicID="smi:org.gfz-potsdam.de/geofon/Magnitude/20210502135732.350764.777031">
<stationCount>31</stationCount>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T13:57:32.350834Z</creationTime>
</creationInfo>
<mag>
<value>4.504973336</value>
<uncertainty>0.2413885735</uncertainty>
</mag>
<type>mb</type>
<originID>smi:org.gfz-potsdam.de/geofon/Origin/20210502135724.073779.776981</originID>
<methodID>smi:org.gfz-potsdam.de/geofon/trimmed_mean</methodID>
</magnitude>
<origin publicID="smi:org.gfz-potsdam.de/geofon/Origin/20210502135724.073779.776981">
<time>
<value>2021-05-02T11:19:04.299081Z</value>
<uncertainty>0.3193747103</uncertainty>
</time>
<longitude>
<value>65.01178741</value>
<uncertainty>3.45298934</uncertainty>
</longitude>
<latitude>
<value>23.98399925</value>
<uncertainty>4.494362354</uncertainty>
</latitude>
<depthType>operator assigned</depthType>
<quality>
<associatedPhaseCount>45</associatedPhaseCount>
<usedPhaseCount>45</usedPhaseCount>
<associatedStationCount>45</associatedStationCount>
<usedStationCount>45</usedStationCount>
<depthPhaseCount>0</depthPhaseCount>
<standardError>0.7319302149</standardError>
<azimuthalGap>85.68946457</azimuthalGap>
<maximumDistance>49.39370728</maximumDistance>
<minimumDistance>8.08137989</minimumDistance>
<medianDistance>37.06093979</medianDistance>
</quality>
<evaluationMode>manual</evaluationMode>
<creationInfo>
<agencyID>GFZ</agencyID>
<creationTime>2021-05-02T13:57:24.074412Z</creationTime>
</creationInfo>
<depth>
<value>10000</value>
<uncertainty>0</uncertainty>
</depth>
<originUncertainty>
<minHorizontalUncertainty>6922.801018</minHorizontalUncertainty>
<maxHorizontalUncertainty>9978.889465</maxHorizontalUncertainty>
<azimuthMaxHorizontalUncertainty>158.6850281</azimuthMaxHorizontalUncertainty>
<preferredDescription>horizontal uncertainty</preferredDescription>
</originUncertainty>
<methodID>smi:org.gfz-potsdam.de/geofon/LOCSAT</methodID>
<earthModelID>smi:org.gfz-potsdam.de/geofon/iasp91</earthModelID>
<evaluationStatus>confirmed</evaluationStatus>
</origin>
<preferredOriginID>smi:org.gfz-potsdam.de/geofon/Origin/20210502135724.073779.776981</preferredOriginID>
<preferredMagnitudeID>smi:org.gfz-potsdam.de/geofon/Magnitude/20210502135732.350764.777031</preferredMagnitudeID>
</event>
</eventParameters>
</q:quakeml>
`
