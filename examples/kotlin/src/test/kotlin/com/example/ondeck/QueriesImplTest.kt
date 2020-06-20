package com.example.ondeck

import com.example.dbtest.DbTestExtension
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.RegisterExtension

class QueriesImplTest {
    companion object {
        @JvmField @RegisterExtension val dbtest = DbTestExtension("src/main/resources/ondeck/schema")
    }

    @Test
    fun testQueries() {
        val q = QueriesImpl(dbtest.getConnection())
        val city = q.createCity(
                slug = "san-francisco",
                name = "San Francisco"
        ).execute()
        val venueId = q.createVenue(
                slug = "the-fillmore",
                name = "The Fillmore",
                city = city.slug,
                spotifyPlaylist = "spotify=uri",
                status = Status.OPEN,
                statuses = listOf(Status.OPEN, Status.CLOSED),
                tags = listOf("rock", "punk")
        ).execute()
        val venue = q.getVenue(
                slug = "the-fillmore",
                city = city.slug
        ).execute()
        assertEquals(venueId, venue.id)

        assertEquals(city, q.getCity(city.slug).execute())
        assertEquals(listOf(VenueCountByCityRow(city.slug, 1)), q.venueCountByCity().execute())
        assertEquals(listOf(city), q.listCities().execute())
        assertEquals(listOf(venue), q.listVenues(city.slug).execute())

        q.updateCityName(slug = city.slug, name = "SF").execute()
        val id = q.updateVenueName(slug = venue.slug, name = "Fillmore").execute()
        assertEquals(venue.id, id)

        q.deleteVenue(venue.slug).execute()
    }
}
