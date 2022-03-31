package zdns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gsmlg-dev/gsmlg-golang/zdns"
	"github.com/libdns/libdns"
)

func (p *Provider) setZdnsToken() {
	zdns.SetToken(p.APIToken)

}

func (p *Provider) getDNSEntries(ctx context.Context, zone string) ([]libdns.Record, error) {

	p.setZdnsToken()

	var records []libdns.Record

	//todo now can only return 100 records
	reqRecords, err := zdns.GetRrByZone(fixZoneName(zone))
	if err != nil {
		return records, err
	}

	for _, entry := range reqRecords {
		record := libdns.Record{
			Name:  entry.Name,
			Value: entry.Rdata,
			Type:  entry.Type,
			TTL:   time.Duration(entry.Ttl) * time.Second,
			ID:    entry.Id,
		}
		records = append(records, record)
	}

	return records, nil
}

func extractRecordName(name string, zone string) string {
	if idx := strings.Index(name, "."+strings.Trim(zone, ".")); idx != -1 {
		return name[:idx]
	}
	return name
}

func (p *Provider) addDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {

	p.setZdnsToken()

	entry := zdns.Rr{
		Name:  extractRecordName(record.Name, fixZoneName(zone)),
		Rdata: record.Value,
		Type:  record.Type,
		Ttl:   int(record.TTL.Seconds()),
	}

	rec, err := zdns.CreateRrInZone(fixZoneName(zone), entry.Name, entry.Type, entry.Ttl, entry.Rdata)
	if err != nil {
		// fmt.Printf("%s, %s, %s, %s, %v", zone, entry.Name, entry.Value, err.Error(), record)
		return record, fmt.Errorf("create record err.Zone:%s, Name: %s, Rdata: %s, Error:%s, %v", zone, entry.Name, entry.Rdata, err.Error(), record)
	}
	record.ID = rec[0].Id

	return record, nil
}

func (p *Provider) removeDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {

	p.setZdnsToken()

	_, err := zdns.DeleteRr(record.ID)
	if err != nil {
		// fmt.Printf("%s, %s, %s, %s, %v", zone, record.Name, record.Value, err.Error(), record)
		return record, fmt.Errorf("remove record err.Zone:%s, Name: %s, Rdata: %s, Error:%s", zone, record.Name, record.Value, err.Error())
	}

	return record, nil
}

func (p *Provider) updateDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {

	p.setZdnsToken()

	entry := zdns.Rr{
		Name:  extractRecordName(record.Name, fixZoneName(zone)),
		Rdata: record.Value,
		Type:  record.Type,
		Ttl:   int(record.TTL.Seconds()),
	}

	_, err := zdns.UpdateRr(fixZoneName(zone), record.ID, entry.Name, entry.Type, entry.Ttl, entry.Rdata)
	if err != nil {
		// fmt.Printf("%s, %s, %s, %s, %v", zone, entry.Name, entry.Value, err.Error(), record)
		return record, fmt.Errorf("update record err.Zone:%s, Name: %s, Rdata: %s, Error:%s, %v", zone, entry.Name, entry.Rdata, err.Error(), record)
	}

	return record, nil
}

func fixZoneName(z string) string {
	return strings.Trim(z, ".") + "."
}
