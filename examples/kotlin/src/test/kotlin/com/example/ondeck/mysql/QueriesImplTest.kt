package com.example.ondeck.mysql

import com.example.dbtest.MysqlDbTestExtension
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension

class QueriesImplTest {
    companion object {
        @JvmField
        @RegisterExtension
        val dbtest = MysqlDbTestExtension("src/main/resources/ondeck/mysql/schema")
    }

    @Test
    fun testQueries() {
        val q = QueriesImpl(dbtest.getConnection())
        q.createCity(
            slug = "san-francisco",
            name = "San Francisco"
        )
        val city = q.listCities()[0]
        val venueId = q.createVenue(
            slug = "the-fillmore",
            name = "The Fillmore",
            city = city.slug,
            spotifyPlaylist = "spotify=uri",
            status = VenueStatus.OPEN,
            statuses = listOf(VenueStatus.OPEN, VenueStatus.CLOSED).joinToString(","),
            tags = listOf("rock", "punk").joinToString(",")
        )
        val venue = q.getVenue(
            slug = "the-fillmore",
            city = city.slug
        )!!
        assertEquals(venueId, venue.id)

        assertEquals(city, q.getCity(city.slug))
        assertEquals(listOf(VenueCountByCityRow(city.slug, 1)), q.venueCountByCity())
        assertEquals(listOf(city), q.listCities())
        assertEquals(listOf(venue), q.listVenues(city.slug))

        q.updateCityName(slug = city.slug, name = "SF")
        q.updateVenueName(slug = venue.slug, name = "Fillmore")
        val fresh = q.getVenue(venue.slug, city.slug)!!
        assertEquals("Fillmore", fresh.name)

        q.deleteVenue(venue.slug, venue.slug)
    }
}
